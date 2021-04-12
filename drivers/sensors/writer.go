package sensors

import (
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/model"
)

type metricWriter struct {
	metric models.Metric
	ctx *Context
}

func (w *metricWriter) Write(value float64) {
	if ch, ok := w.ctx.Pipe[w.metric]; ok {
		ch <- model.MetricReading {
			Source: w.ctx.SensorID,
			Value: value,
		}
	}
}

func (w *metricWriter) Write32(value float32) {
	w.Write(float64(value))
}

func (w *metricWriter) WriteWithError(value float64, err error) {
	if err != nil {
		w.ctx.Error(err)
		return
	}
	w.Write(value)
}

func (w *metricWriter) Write32WithError(value float32, err error) {
	w.WriteWithError(float64(value), err)
}
