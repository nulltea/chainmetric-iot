package sensor

import (
	"context"

	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-iot/shared"
)

// Context defines structure for the sensor.Sensor reading context.
type Context struct {
	context.Context
	SensorID string
	Pipe     ReadingsPipe
}

// NewReaderContext constructs new Context instance based on given `parent` context for the given sensor.Sensor.
func NewReaderContext(parent context.Context, sensor Sensor) *Context {
	return &Context{
		Context: parent,
		SensorID: sensor.ID(),
		Pipe: make(ReadingsPipe),
	}
}

// WriterFor returns MetricWriter for a given models.Metric.
func (c *Context) WriterFor(metric models.Metric) *MetricWriter {
	return &MetricWriter{
		metric,
		c,
	}
}

// Error wraps `err` logging with sensor.Sensor metadata.
func (c *Context) Error(err error) {
	if err != nil {
		shared.Logger.Errorf("%v: %v", c.SensorID, err)
	}
}

// Warning wraps `msg` logging with sensor.Sensor metadata.
func (c *Context) Warning(msg string) {
	shared.Logger.Errorf("%v: %v", c.SensorID, msg)
}

// Info wraps `info` logging with sensor.Sensor metadata.
func (c *Context) Info(info string) {
	shared.Logger.Infof("%v: %v", c.SensorID, info)
}
