package sensor

import (
	"fmt"

	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/model"
)

type metricWriter struct {
	metric models.Metric
	ctx *Context
}

func (w *metricWriter) Write(v interface{}) {
	var value float64

	switch t := v.(type) {
	case float64:
		value = t
	case float32:
		value = float64(t)
	case int:
		value = float64(t)
	case int32:
		value = float64(t)
	case int64:
		value = float64(t)
	case uint8:
		value = float64(t)
	case uint16:
		value = float64(t)
	default:
		w.ctx.Error(fmt.Errorf("value type is not supported: %T", t))
		return
	}


	if ch, ok := w.ctx.Pipe[w.metric]; ok {
		ch <- model.SensorReading{
			Source: w.ctx.SensorID,
			Value:  value,
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
