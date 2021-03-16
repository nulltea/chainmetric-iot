package sensors

import (
	"context"
	"time"

	"github.com/op/go-logging"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/model"
)

type Context struct {
	Parent    context.Context
	SensorID  string
	Logger    *logging.Logger
	Pipe      model.MetricReadingsPipe
}

func (c *Context) For(metric models.Metric) *metricWriter {
	return &metricWriter{
		metric,
		c,
	}
}

func (c *Context) Error(err error) {
	if err != nil {
		c.Logger.Errorf("%v: %v", c.SensorID, err)
	}
}

func (c *Context) Info(info string) {
	c.Logger.Infof("%v: %v", c.SensorID, info)
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
