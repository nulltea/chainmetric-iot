package readings

import (
	"context"
	"sync"
	"time"

	"github.com/op/go-logging"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/config"
	"github.com/timoth-y/iot-blockchain-sensorsys/model"
	"github.com/timoth-y/iot-blockchain-sensorsys/readings/receiver"
	"github.com/timoth-y/iot-blockchain-sensorsys/readings/sensor"
	"github.com/timoth-y/iot-blockchain-sensorsys/sensors"
)

type Context struct {
	Parent context.Context
	Logger *logging.Logger
	WaitGroup *sync.WaitGroup
	Config config.Config
}

func NewContext(parent context.Context) *Context {
	return &Context{
		Parent: parent,
	}
}

func (c *Context) ForSensor(s sensors.Sensor) *sensor.Context {
	return &sensor.Context {
		Parent: c,
		Logger: c.Logger,
		SensorID: s.ID(),
		Pipe: make(model.MetricReadingsPipe),
	}
}

func (c *Context) ForRequest(metrics []models.Metric) *receiver.Context {
	pipe := make(model.MetricReadingsPipe)
	for _, metric := range metrics {
		pipe[metric] = make(chan model.MetricReading, 3)
	}
	return &receiver.Context {
		Parent: c,
		Logger: c.Logger,
		Pipe: pipe,
	}
}

func (c *Context) SetLogger(logger *logging.Logger) *Context {
	c.Logger = logger
	return c
}

func (c *Context) SetConfig(configPath string) *Context {
	cnf, err := config.ReadConfig(configPath); if err != nil {
		c.Fatal(err)
	}
	c.Config = cnf
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
