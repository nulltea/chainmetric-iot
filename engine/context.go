package engine

import (
	"context"
	"sync"
	"time"

	"github.com/op/go-logging"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/sensors"
	"github.com/timoth-y/iot-blockchain-sensorsys/model"
)

type Context struct {
	Parent context.Context
	Logger *logging.Logger
	WaitGroup *sync.WaitGroup
}

func NewContext(parent context.Context) *Context {
	return &Context{
		Parent: parent,
	}
}

func (c *Context) ForSensor(s sensors.Sensor) *sensors.Context {
	return &sensors.Context{
		Parent: c,
		Logger: c.Logger,
		SensorID: s.ID(),
		Pipe: make(model.MetricReadingsPipe),
	}
}

func (c *Context) SetLogger(logger *logging.Logger) *Context {
	c.Logger = logger
	return c
}

func (c *Context) Error(err error) {
	if err != nil {
		c.Logger.Errorf("worker: %v", err)
	}
}

func (c *Context) Fatal(err error) {
	if err != nil {
		c.Logger.Fatalf("worker: %v", err)
	}
}

func (c *Context) Info(info string) {
	c.Logger.Infof("worker: %v", info)
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
