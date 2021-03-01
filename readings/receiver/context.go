package receiver

import (
	"sensorsys/model"
	"sensorsys/readings"
	"sensorsys/readings/sensor"
)

type Context struct {
	*readings.Context
	Pipe   model.MetricReadingsPipe
}

func (c *Context) ForSensor(s sensor.Sensor) *sensor.Context {
	return &sensor.Context{
		Context: c.Context,
		Pipe: c.Pipe,
	}
}
