package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/device"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/display"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/periphery"
	"github.com/timoth-y/iot-blockchain-sensorsys/engine"
	"github.com/timoth-y/iot-blockchain-sensorsys/gateway/blockchain"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/config"
	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

var (
	Client = blockchain.NewBlockchainClient()
	Display = display.NewST7789()
	Reader = engine.NewSensorsReader()
	Context = engine.NewContext(context.Background()).
		SetLogger(shared.Logger)
	Device = device.NewDevice().
		SetClient(Client).
		SetDisplay(Display).
		SetReader(Reader)
)

func init() {
	shared.InitLogger()
  shared.InitConfig()
	periphery.Init()
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
	var bc config.BlockchainConfig
	if err := viper.UnmarshalKey("gateway", &bc); err != nil {
		shared.Logger.Fatal(errors.Wrap(err, "failed parse blockchain config"))
	}
	if err := Client.Init(bc); err != nil {
		shared.Logger.Fatal(errors.Wrap(err, "failed initializing blockchain client"))
	}

	var dc config.DisplayConfig
	if err := viper.UnmarshalKey("display", &dc); err != nil {
		shared.Logger.Fatal(errors.Wrap(err, "failed parse display config"))
	}
	if err := Display.Init(dc); err != nil {
		shared.Logger.Fatal(errors.Wrap(err, "failed initializing display"))
	}

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

