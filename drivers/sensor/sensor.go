package sensor

import (
	"github.com/timoth-y/chainmetric-core/models"
)

type Sensor interface {
	ID() string
	Init() error
	Harvest(ctx *Context)
	Metrics() []models.Metric
	Verify() bool
	Active() bool
	Close() error
}

