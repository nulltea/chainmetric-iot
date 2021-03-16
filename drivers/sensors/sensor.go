package sensors

import (
	"github.com/timoth-y/iot-blockchain-contracts/models"
)

type Sensor interface {
	ID() string
	Init() error
	Harvest(ctx *Context)
	Metrics() []models.Metric
	Active() bool
	Close() error
}
