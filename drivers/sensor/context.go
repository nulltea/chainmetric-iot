package sensor

import (
	"context"
	"time"

	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-sensorsys/model"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

type Context struct {
	Parent    context.Context
	SensorID  string
	Pipe      model.SensorReadingsPipe
}

func (c *Context) For(metric models.Metric) *metricWriter {
	return &metricWriter{
		metric,
		c,
	}
}

func (c *Context) Error(err error) {
	if err != nil {
		shared.Logger.Errorf("%v: %v", c.SensorID, err)
	}
}

func (c *Context) Warning(msg string) {
	shared.Logger.Errorf("%v: %v", c.SensorID, msg)
}

func (c *Context) Info(info string) {
	shared.Logger.Infof("%v: %v", c.SensorID, info)
}

func (c *Context) SetTimeout(timeout time.Duration) (*Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.Parent, timeout)
	c.Parent = ctx
	return c, cancel
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
