package sensor

import (
	"context"
	"sync"
	"time"

	"github.com/op/go-logging"

	"sensorsys/model"
)

type Context struct {
	Parent    context.Context
	SensorID  string
	Logger    *logging.Logger
	WaitGroup *sync.WaitGroup
	Pipe      model.MetricReadingsPipe
}

func (c *Context) For(metric model.Metric) *metricWriter {
	return &metricWriter {
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
