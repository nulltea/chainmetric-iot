package sensor

import (
	"sensorsys/model"
	"sensorsys/readings"
)

type Context struct {
	*readings.Context
	Pipe   model.MetricReadingsPipe
}

func (c *Context) For(metric model.Metric) *metricWriter {
	return &metricWriter {
		metric,
		c,
	}
}
