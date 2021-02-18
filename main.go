package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/op/go-logging"
	"github.com/stianeikeland/go-rpio/v4"

	"sensor/readings"
)

var (
	logger = logging.MustGetLogger("sensor")
)

func main() {
	done := make(chan struct{}, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go process()
	go shutdown(quit, done)

	<-done
	logger.Info("Shutdown")
}

func process() {
	reader := readings.NewSensorsReader(context.Background())
	reader.SubscribeToTemperatureReadings("DHT11", 4)

	for {
		readings := reader.Execute()
		fmt.Println(readings)
		time.Sleep(5 * time.Second)
	}
}

func shutdown(quit chan os.Signal, done chan struct{}) {
	<-quit
	logger.Info("Shutting down...")

	err := rpio.Close(); if err != nil {
		logger.Error(err)
	}

	close(done)
}