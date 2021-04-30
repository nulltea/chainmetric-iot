package main

import (
	"os"
	"os/signal"

	"github.com/spf13/viper"

	dev "github.com/timoth-y/chainmetric-sensorsys/drivers/device"
	displays "github.com/timoth-y/chainmetric-sensorsys/drivers/display"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/periphery"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensors"
	"github.com/timoth-y/chainmetric-sensorsys/engine"
	"github.com/timoth-y/chainmetric-sensorsys/gateway/blockchain"
	"github.com/timoth-y/chainmetric-sensorsys/model/config"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

var (
	bcf config.BlockchainConfig
	dcf config.DisplayConfig

	client  *blockchain.Client
	reader  *engine.SensorsReader
	display displays.Display
	device  *dev.Device
)

func init() {
	shared.InitCore()
	periphery.Init()

	shared.MustUnmarshalFromConfig("gateway", &bcf)
	shared.MustUnmarshalFromConfig("display", &dcf)

	client = blockchain.NewClient(bcf)
	reader = engine.NewSensorsReader()
	display = displays.NewEInk(dcf)
	device = dev.New().
		SetClient(client).
		SetReader(reader).
		SetDisplay(display)
}

func main() {
	done := make(chan struct{}, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go startup()
	go shutdown(quit, done)

	<-done
	shared.Logger.Info("Bye!")
}

func startup() {
	if dcf.Enabled {
		shared.MustExecute(display.Init, "failed initializing display")
	}

	if viper.GetBool("mocks.debug_env") {
		device.RegisterStaticSensors(sensors.NewStaticSensorMock())
	}

	shared.MustExecute(client.Init, "failed initializing blockchain client")
	shared.MustExecute(device.Init, "failed to initialize device")
	shared.MustExecute(device.CacheBlockchainState, "failed to cache the state of blockchain")

	device.WatchForBlockchainEvents()
	device.Operate()
}

func shutdown(quit chan os.Signal, done chan struct{}) {
	<-quit
	shared.Logger.Info("Shutting down...")

	if dcf.Enabled {
		shared.Execute(display.ClearAndRefresh, "error during clearing display")
		shared.Execute(display.Close, "error during closing connection to display")
	}

	shared.Execute(device.NotifyOff, "error during emitting 'off' event")
	shared.Execute(device.Close, "error during device shutdown")

	client.Close()
	shared.CloseCore()

	close(done)
}
