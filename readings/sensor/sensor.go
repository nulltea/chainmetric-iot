package sensor

import (
	"sensorsys/model"
)

type Sensor interface {
	ID() string
	Init() error
	Harvest(ctx *Context)
	Metrics() []model.Metric
	Active() bool
	Close() error
}
