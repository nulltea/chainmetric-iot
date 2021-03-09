package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"time"

	"github.com/d2r2/go-logger"
	"github.com/op/go-logging"
	"periph.io/x/periph/host"

	"github.com/timoth-y/iot-blockchain-sensorsys/mocks"
	"github.com/timoth-y/iot-blockchain-sensorsys/model"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
	"github.com/timoth-y/iot-blockchain-sensorsys/readings"
	"github.com/timoth-y/iot-blockchain-sensorsys/sensors"
	"github.com/timoth-y/iot-blockchain-sensorsys/utils"
)

const (
	format = "%{color}%{time:2006.01.02 15:04:05} %{id:04x} %{level}%{color:reset} [%{module}] %{color:bold}%{shortfunc}%{color:reset} -> %{message}"
)

var (
	Logger = logging.MustGetLogger("sensorsys")
	ctx = readings.NewContext(context.Background()).
		SetLogger(Logger).
		SetConfig("config.yaml")
	reader = readings.NewSensorsReader(ctx)
)

func init() {
	initLogging()

	if _, err := host.Init(); err != nil {
		Logger.Fatal(err)
	}
}

func main() {
	done := make(chan struct{}, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	utils.GenerateDeviceSignatureInQR()

	go run()
	go shutdown(quit, done)

	<-done
	Logger.Info("Shutdown")
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
		Logger.Info(string(s))
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
		Logger.Info(string(s))
	}, 2 * time.Second,
		metrics.Temperature,
		metrics.Humidity,
		metrics.Luminosity,
	)
}

func shutdown(quit chan os.Signal, done chan struct{}) {
	<-quit
	Logger.Info("Shutting down...")

	reader.Clean()

	close(done)
}

func initLogging() {
	backend := logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0),
		logging.MustStringFormatter(format))

	logging.SetBackend(backend)

	level, err := logging.LogLevel(os.Getenv("LOGGING")); if err != nil {
		level = logging.INFO
	}
	logging.SetLevel(level, "sensorsys")

	logger.ChangePackageLogLevel("dht", logger.ErrorLevel)
	logger.ChangePackageLogLevel("i2c", logger.ErrorLevel)
}
