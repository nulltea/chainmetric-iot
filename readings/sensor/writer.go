package sensor

import (
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/model"
)

type metricWriter struct {
	metric models.Metric
	ctx *Context
}

func (w *metricWriter) Write(value interface{}) {
	if ch, ok := w.ctx.Pipe[w.metric]; ok {
		ch <- model.MetricReading {
			Source: w.ctx.SensorID,
			Value: value,
		}
	}
}

func (w *metricWriter) WriteWithError(value interface{}, err error) {
	if err != nil {
		w.ctx.Error(err)
		return
	}
	w.Write(value)
}

