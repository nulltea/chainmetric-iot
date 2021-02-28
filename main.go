package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"

	log "github.com/d2r2/go-logger"
	"github.com/op/go-logging"

	"sensorsys/model"
	"sensorsys/model/metrics"
	"sensorsys/sensors"
	"sensorsys/worker"
)

var (
	logger = logging.MustGetLogger("sensor")
	ctx = worker.NewContext(context.Background()).SetLogger(logger)
	reader = worker.NewSensorsReader(ctx)
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

	reader.SubscribeReceiver(func(readings model.MetricReadings) {
		s, _ := json.MarshalIndent(readings, "", "\t")
		logger.Info(string(s))
	}, metrics.Temperature,
		metrics.Humidity,
		metrics.Luminosity,
		metrics.UVLight,
		metrics.VisibleLight,
		metrics.IRLight,
		metrics.HeartRate,
		metrics.BloodOxidation,
		metrics.AirCO2Concentration,
		metrics.AirTVOCsConcentration,
	)
}

func shutdown(quit chan os.Signal, done chan struct{}) {
	<-quit
	logger.Info("Shutting down...")

	reader.Clean()

	close(done)
}
