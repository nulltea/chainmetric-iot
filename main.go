package main

import (
	"os"
	"os/signal"

	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-iot/controllers/device/modules"
	"github.com/timoth-y/chainmetric-iot/controllers/gui"
	core "github.com/timoth-y/chainmetric-iot/core/dev"
	dsp "github.com/timoth-y/chainmetric-iot/drivers/display"
	"github.com/timoth-y/chainmetric-iot/network/localnet"

	dev "github.com/timoth-y/chainmetric-iot/controllers/device"
	"github.com/timoth-y/chainmetric-iot/drivers/sensors"
	"github.com/timoth-y/chainmetric-iot/model/config"
	"github.com/timoth-y/chainmetric-iot/network/blockchain"
	"github.com/timoth-y/chainmetric-iot/shared"
)

var (
	dcf config.DisplayConfig

	display core.Display
	device  *dev.Device

	done = make(chan struct{}, 1)
	quit = make(chan os.Signal, 1)
)

func init() {
	shared.InitCore()

	shared.MustUnmarshalFromConfig("display", &dcf)

	device = dev.New(
		modules.WithLifecycleManager(),
		modules.WithEngineOperator(),
		modules.WithCacheManager(),
		modules.WithEventsObserver(),
		modules.WithHotswapDetector(),
		modules.WithRemoteController(),
		modules.WithLocationManager(),
		modules.WithPowerManager(),
		modules.WithFailoverHandler(),
	)

	display = dsp.NewEInk(dcf)
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
		return blockchain.Init()
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
