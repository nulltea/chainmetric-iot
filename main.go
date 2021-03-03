package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"time"

	log "github.com/d2r2/go-logger"
	"github.com/op/go-logging"

	"sensorsys/mocks"
	"sensorsys/model"
	"sensorsys/model/metrics"
	"sensorsys/readings"
	"sensorsys/sensors"
)

var (
	logger = logging.MustGetLogger("sensor")
	ctx = readings.NewContext(context.Background()).
		SetLogger(logger)
	reader = readings.NewSensorsReader(ctx)
)

func main() {
	done := make(chan struct{}, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	log.ChangePackageLogLevel("dht", log.ErrorLevel)
	log.ChangePackageLogLevel("i2c", log.ErrorLevel)

	go run()
	go shutdown(quit, done)

	<-done
	logger.Info("Shutdown")
}

func run() {
	reader.RegisterSensors(
		sensors.NewDHT22(5),
		sensors.NewMAX44009(0x4A, 1),
		sensors.NewMAX30102(0x57, 2),
		sensors.NewCCS811(0x5a, 3),
		sensors.NewSI1145(0x60, 4),
	)

	go reader.Process()

	reader.SubscribeReceiver(func(readings model.MetricReadings) {
		s, _ := json.MarshalIndent(readings, "", "\t")
		logger.Info(string(s))
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
		logger.Info(string(s))
	}, 2 * time.Second,
		metrics.Temperature,
		metrics.Humidity,
		metrics.Luminosity,
	)
}

func shutdown(quit chan os.Signal, done chan struct{}) {
	<-quit
	logger.Info("Shutting down...")

	reader.Clean()

	close(done)
}
