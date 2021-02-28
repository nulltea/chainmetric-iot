package worker

import (
	"context"
	"time"

	"github.com/op/go-logging"

	"sensorsys/config"
	"sensorsys/model"
)

type Context struct {
	Parent   context.Context
	SensorID string
	Logger   *logging.Logger
	Pipe     model.MetricReadingsPipe
	Config   config.Config
}

func NewContext(parent context.Context) *Context {
	return &Context{
		Parent: parent,
	}
}

func (c *Context) ForSensor(sensor Sensor) *Context {
	return &Context {
		Parent: c,
		SensorID: sensor.ID(),
	}
}

func (c *Context) SetLogger(logger *logging.Logger) *Context {
	c.Logger = logger
	return c
}

func (c *Context) SetConfig(config config.Config) *Context {
	c.Config = config
	return c
}

func (c *Context) Error(err error) {
	if err != nil {
		c.Logger.Errorf("%v: %v", c.SensorID, err)
	}
}

func (c *Context) Info(info string) {
	c.Logger.Infof("%v: %v", c.SensorID, info)
}

func (c *Context) For(metric model.Metric) *metricWriter {
	return &metricWriter {
		metric,
		c,
	}
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.Done()
}

func (c *Context) Err() error {
	return c.Err()
}

func (c *Context) Value(key interface{}) interface{} {
	return c.Value(key)
}


