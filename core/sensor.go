package core

import (
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-sensorsys/core/sensor"
)

// Sensor defines base methods for controlling sensor device.
type Sensor interface {
	// ID returns unique Sensor identifier key.
	ID() string
	// Init performers initialization sequence of the Sensor device.
	Init() error
	// Harvest collects all available models.Metric from Sensor device and dumps them into the context.
	Harvest(ctx *sensor.Context)
	// Metrics return all available models.Metric to reads from Sensor device.
	Metrics() []models.Metric
	// Verify checks whether the driver is compatible with Sensor device.
	Verify() bool
	// Active checks whether the Sensor device is connected and active.
	Active() bool
	// Close closes connection to Sensor device and clears allocated resources.
	Close() error
}

