package main

import (
	"os"
	"os/signal"

	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/device/modules"
	"github.com/timoth-y/chainmetric-sensorsys/network/localnet"

	dev "github.com/timoth-y/chainmetric-sensorsys/drivers/device"
	dsp "github.com/timoth-y/chainmetric-sensorsys/drivers/display"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/gui"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensors"
	"github.com/timoth-y/chainmetric-sensorsys/model/config"
	"github.com/timoth-y/chainmetric-sensorsys/network/blockchain"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

var (
	bcf config.BlockchainConfig
	dcf config.DisplayConfig

	display  dsp.Display
	device   *dev.Device

	done = make(chan struct{}, 1)
	quit = make(chan os.Signal, 1)
)

func init() {
	shared.InitCore()

	shared.MustUnmarshalFromConfig("blockchain", &bcf)
	shared.MustUnmarshalFromConfig("display", &dcf)

	display = dsp.NewEInk(dcf)
	device = dev.New(
		modules.WithLifecycleManager(),
		modules.WithEngineOperator(),
		modules.WithCacheManager(),
		modules.WithEventsObserver(),
		modules.WithHotswapDetector(),
		modules.WithRemoteCommandsHandler(),
		modules.WithLocationManager(),
		modules.WithPowerManager(),
		modules.WithFailoverHandler(),
	)

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

	shared.MustExecute(func() error {
		return blockchain.Init(bcf)
	}, "failed initializing blockchain client")

	device.Start()
}

func shutdown() {
	<-quit
	shared.Logger.Info("Shutting down...")

	if dcf.Enabled {
		shared.Execute(display.ClearAndRefresh, "error during clearing display")
		shared.Execute(display.Close, "error during closing connection to display")
	}

	shared.Execute(localnet.Close, "error during closing local network")
	shared.Execute(device.Close, "error during device shutdown")

	blockchain.Close()
	shared.CloseCore()

	close(done)
}
