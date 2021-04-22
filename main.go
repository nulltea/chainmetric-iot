package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/pkg/errors"

	"github.com/timoth-y/iot-blockchain-sensorsys/config"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/device"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/display"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/peripherals"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/sensors"
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
	ADC *peripherals.AnalogMCP3208
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

	assignAnalogSensors()

	if err := Reader.Init(Context); err != nil {
		shared.Logger.Fatal(errors.Wrap(err, "failed initializing reader engine"))
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

func assignAnalogSensors() {
	ADC = peripherals.NewAnalogMCP3208("SPI0.0", 25)

	if err := ADC.Init(); err != nil {
		shared.Logger.Fatal(err, "failed to init MCP3208")
	}

	Device.RegisterStaticSensors(sensors.NewAnalogPZT(ADC.GetChannel(0)))
	Device.RegisterStaticSensors(sensors.NewAnalogHall(ADC.GetChannel(7)))
}

func shutdown(quit chan os.Signal, done chan struct{}) {
	<-quit
	shared.Logger.Info("Shutting down...")

	if Device != nil {
		if err := Device.Off(); err != nil {
			shared.Logger.Error(err)
		}

		if err := Device.Close(); err != nil {
			shared.Logger.Error(err)
		}
	}


	Client.Close()

	close(done)
}

