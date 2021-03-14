package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"time"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/device"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/display"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/sensors"
	"github.com/timoth-y/iot-blockchain-sensorsys/engine"
	"github.com/timoth-y/iot-blockchain-sensorsys/gateway/blockchain"
	"github.com/timoth-y/iot-blockchain-sensorsys/mocks"
	"github.com/timoth-y/iot-blockchain-sensorsys/model"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

var (
	Device *device.Device
	ctx = engine.NewContext(context.Background()).
		SetLogger(shared.Logger).
		SetConfig("config.yaml")
	reader = engine.NewSensorsReader(ctx)
)

func init() {
	shared.InitLogger()
	shared.InitPeriphery()
}

func main() {
	done := make(chan struct{}, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go run()
	go shutdown(quit, done)

	<-done
	shared.Logger.Info("Shutdown")
}

func run() {
	Device = device.NewDevice()

	Device.SetConfig(ctx.Config)

	client := blockchain.NewBlockchainClient()
	if err := client.Init(ctx.Config.Gateway); err != nil {
		shared.Logger.Fatal(err)
	}
	shared.Logger.Debug("init blockchain client ended")

	Device.SetClient(client)

	dp := display.NewST7789()

	if err := dp.Init(ctx.Config.Display); err != nil {
		shared.Logger.Fatal(err)
	}

	Device.SetDisplay(dp)

	if err := Device.Init(); err != nil {
		shared.Logger.Fatal(err)
	}

	Device.RegisterStaticSensors(sensors.NewDHT22(ctx.Config.Sensors.DHT22.Pin))
	reader.RegisterSensors(Device.SupportedSensors()...)

	go reader.Process()

	reader.SubscribeReceiver(func(readings model.MetricReadings) {
		s, _ := json.MarshalIndent(readings, "", "\t")
		shared.Logger.Info(string(s))
	}, 3 * time.Second,
		metrics.Temperature,
		metrics.Humidity,
		metrics.Luminosity,
		metrics.UVLight,
		metrics.VisibleLight,
		metrics.IRLight,
	)
}

func mock() {
	reader.RegisterSensors(
		mocks.NewMockSensor(500 * time.Millisecond, metrics.Luminosity),
		mocks.NewMockSensor(800 * time.Millisecond, metrics.Humidity),
		mocks.NewMockSensor(1000 * time.Millisecond, metrics.Temperature),
	)

	go reader.Process()

	reader.SubscribeReceiver(func(readings model.MetricReadings) {
		s, _ := json.MarshalIndent(readings, "", "\t")
		shared.Logger.Info(string(s))
	}, 2 * time.Second,
		metrics.Temperature,
		metrics.Humidity,
		metrics.Luminosity,
	)
}

func shutdown(quit chan os.Signal, done chan struct{}) {
	<-quit
	shared.Logger.Info("Shutting down...")

	reader.Clean()

	Device.Close()

	close(done)
}

