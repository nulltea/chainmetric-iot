package engine

import (
	"context"
	"sync"
	"time"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
	"github.com/timoth-y/chainmetric-sensorsys/model"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

type Context struct {
	Parent context.Context
	WaitGroup *sync.WaitGroup
}

func NewContext(parent context.Context) *Context {
	return &Context{
		Parent: parent,
	}
}

func (c *Context) ForSensor(s sensor.Sensor) *sensor.Context {
	return &sensor.Context{
		Parent: c,
		SensorID: s.ID(),
		Pipe: make(model.SensorReadingsPipe),
	}
}

func (c *Context) Error(err error) {
	if err != nil {
		shared.Logger.Errorf("worker: %v", err)
	}
}

func (c *Context) Fatal(err error) {
	if err != nil {
		shared.Logger.Fatalf("worker: %v", err)
	}
}

func (c *Context) Info(info string) {
	shared.Logger.Infof("worker: %v", info)
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
