package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	log "github.com/d2r2/go-logger"
	"github.com/op/go-logging"
	"github.com/stianeikeland/go-rpio/v4"

	"sensorsys/readings"
)

var (
	logger = logging.MustGetLogger("sensor")
	reader = readings.NewSensorsReader(context.Background())
)


func main() {
	done := make(chan struct{}, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	log.ChangePackageLogLevel("dht", log.ErrorLevel)
	log.ChangePackageLogLevel("i2c", log.ErrorLevel)

	go process()
	go shutdown(quit, done)

	<-done
	logger.Info("Shutdown")
}

func process() {
	if err := reader.SubscribeToTemperatureReadings("DHT22", 5); err != nil {
		logger.Error(err)
	}
	if err := reader.SubscribeToLuminosityReadings(0x4A, 1); err != nil {
		logger.Error(err)
	}
	if err := reader.SubscribeToAirQualityReadings(0x5a, 3); err != nil {
		logger.Error(err)
	}
	if err := reader.SubscribeToAmbientLightReadings(0x60, 4); err != nil {
		logger.Error(err)
	}
	if err := reader.SubscribeToVitalsReadings(0x57, 2); err != nil {
		logger.Error(err)
	}
	// reader.SubscribeToPressureReadings("BMP280", 0x76, 3)
	err := rpio.Open(); if err != nil{
		logger.Fatal(err)
	}
	pin := rpio.Pin(4)
	pin.Input()
	for {
		readings := reader.Execute()
		// readings := pin.Read()
		fmt.Println(readings)
		time.Sleep(1 * time.Second)
	}
}

func shutdown(quit chan os.Signal, done chan struct{}) {
	<-quit
	logger.Info("Shutting down...")

	rpio.Close()
	reader.Clean()

	close(done)
}
