package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/device"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/display"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/sensors"
	"github.com/timoth-y/iot-blockchain-sensorsys/engine"
	"github.com/timoth-y/iot-blockchain-sensorsys/gateway/blockchain"
	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

var (
	Device = device.NewDevice().
		SetConfig(ctx.Config).
		SetClient(blockchain.NewBlockchainClient()).
		SetDisplay(display.NewST7789()).
		RegisterStaticSensors(sensors.NewDHT22(ctx.Config.Sensors.DHT22.Pin))

	ctx = engine.NewContext(context.Background()).
		SetLogger(shared.Logger).
		SetConfig("config.yaml")
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
	Device.Init()
}

func shutdown(quit chan os.Signal, done chan struct{}) {
	<-quit
	shared.Logger.Info("Shutting down...")

	if err := Device.Close(); err != nil {
		shared.Logger.Error(err)
	}

	close(done)
}

