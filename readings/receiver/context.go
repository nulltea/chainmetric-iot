package receiver

import (
	"context"
	"sync"
	"time"

	"github.com/op/go-logging"

	"github.com/timoth-y/iot-blockchain-sensorsys/model"
	"github.com/timoth-y/iot-blockchain-sensorsys/readings/sensor"
	"github.com/timoth-y/iot-blockchain-sensorsys/sensors"
)

type Context struct {
	Parent    context.Context
	SensorID  string
	Logger    *logging.Logger
	WaitGroup *sync.WaitGroup
	Pipe      model.MetricReadingsPipe
}

func (c *Context) ForSensor(s sensors.Sensor) *sensor.Context {
	return &sensor.Context {
		Parent: c,
		Pipe: c.Pipe,
		SensorID: s.ID(),
		Logger: c.Logger,
	}
}

func (c *Context) Error(err error) {
	if err != nil {
		c.Logger.Error(err)
	}
}

func (c *Context) Info(info string) {
	c.Logger.Info(info)
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.Parent.Deadline()
}

func (c *Context) Done() <- chan struct{} {
	return c.Parent.Done()
}

func (c *Context) Err() error {
	return c.Parent.Err()
}

func (c *Context) Value(key interface{}) interface{} {
	return c.Parent.Value(key)
}


