package sensors

import (
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/readings/sensor"
)

type Sensor interface {
	ID() string
	Init() error
	Harvest(ctx *sensor.Context)
	Metrics() []models.Metric
	Active() bool
	Close() error
}
