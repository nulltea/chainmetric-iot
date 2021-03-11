package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"time"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/sensors"
	"github.com/timoth-y/iot-blockchain-sensorsys/engine"
	"github.com/timoth-y/iot-blockchain-sensorsys/mocks"
	"github.com/timoth-y/iot-blockchain-sensorsys/model"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

var (
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
	reader.RegisterSensors(
		sensors.NewDHT22(5),
		sensors.NewMAX44009(0x4A, 1),
		sensors.NewMAX30102(0x57, 2),
		sensors.NewCCS811(0x5A, 3),
		sensors.NewSI1145(0x60, 4),
	)

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

	close(done)
}

