package main

import (
	"os"
	"os/signal"

	"github.com/spf13/viper"

	dev "github.com/timoth-y/chainmetric-sensorsys/drivers/device"
	dsp "github.com/timoth-y/chainmetric-sensorsys/drivers/display"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/gui"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensors"
	"github.com/timoth-y/chainmetric-sensorsys/engine"
	"github.com/timoth-y/chainmetric-sensorsys/model/config"
	"github.com/timoth-y/chainmetric-sensorsys/network/blockchain"
	"github.com/timoth-y/chainmetric-sensorsys/network/local"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

var (
	bcf config.BlockchainConfig
	dcf config.DisplayConfig

	display  dsp.Display
	client   *blockchain.Client
	reader   *engine.SensorsReader
	localnet *local.Client
	device   *dev.Device

	done = make(chan struct{}, 1)
	quit = make(chan os.Signal, 1)
)

func init() {
	shared.InitCore()

	shared.MustUnmarshalFromConfig("blockchain", &bcf)
	shared.MustUnmarshalFromConfig("display", &dcf)

	client = blockchain.NewClient(bcf)
	reader = engine.NewSensorsReader()
	display = dsp.NewEInk(dcf)
	localnet = local.NewClient()
	device = dev.New().
		SetClient(client).
		SetEngine(reader).
		SetLocalNet(localnet)

	gui.Init(display)
}

func main() {
	signal.Notify(quit, os.Interrupt)

	go startup()
	go shutdown()

	<-done
	shared.Logger.Info("Goodbye.")
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
	shared.MustExecute(device.ListenRemoteCommands, "failed to start remote commands listener")

	device.WatchForBlockchainEvents()
	device.Operate()
}

func shutdown() {
	<-quit
	shared.Logger.Info("Shutting down...")

	if dcf.Enabled {
		shared.Execute(display.ClearAndRefresh, "error during clearing display")
		shared.Execute(display.Close, "error during closing connection to display")
	}

	shared.Execute(device.Close, "error during closing local network")

	shared.Execute(device.NotifyOff, "error during emitting 'off' event")
	shared.Execute(device.Close, "error during device shutdown")

	client.Close()
	shared.CloseCore()

	close(done)
}
