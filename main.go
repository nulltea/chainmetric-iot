package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/pkg/errors"

	"github.com/timoth-y/iot-blockchain-sensorsys/config"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/device"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/display"
	"github.com/timoth-y/iot-blockchain-sensorsys/engine"
	"github.com/timoth-y/iot-blockchain-sensorsys/gateway/blockchain"
	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

var (
	Config = config.MustReadConfig("config.yaml")
	Client = blockchain.NewBlockchainClient()
	Display = display.NewST7789()
	Reader = engine.NewSensorsReader()
	Context = engine.NewContext(context.Background()).
		SetConfig(Config).
		SetLogger(shared.Logger)
	Device = device.NewDevice().
		SetConfig(Config).
		SetClient(Client).
		SetDisplay(Display).
		SetReader(Reader)
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
	if err := Client.Init(Config.Gateway); err != nil {
		shared.Logger.Fatal(errors.Wrap(err, "failed initializing blockchain client"))
	}

	if err := Display.Init(Config.Display); err != nil {
		shared.Logger.Fatal(errors.Wrap(err, "failed initializing display"))
	}

	if err := Reader.Init(Context); err != nil {
		shared.Logger.Fatal(err)
	}

	if err := Device.Init(); err != nil {
		shared.Logger.Fatal(errors.Wrap(err, "failed to initialize device"))
	}

	if err := Device.CacheBlockchainState(); err != nil {
		shared.Logger.Fatal(errors.Wrap(err, "failed to cache the state of blockchain"))
	}

	Device.WatchForBlockchainEvents()

	Device.Operate()
}

func shutdown(quit chan os.Signal, done chan struct{}) {
	<-quit
	shared.Logger.Info("Shutting down...")

	if err := Device.Close(); err != nil {
		shared.Logger.Error(err)
	}

	if err := Device.Off(); err != nil {
		shared.Logger.Error(err)
	}

	Client.Close()

	close(done)
}

