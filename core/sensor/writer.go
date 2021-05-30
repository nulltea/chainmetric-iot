package sensor

import (
	"fmt"

	"github.com/timoth-y/chainmetric-core/models"
)

// MetricWriter defines object capable of dumping reading results from core.Sensor
// to a specific ReadingsPipe for models.Metric data.
type MetricWriter struct {
	metric models.Metric
	ctx *Context
}

// Write writes reading results from core.Sensor with required type conversation.
func (w *MetricWriter) Write(v interface{}) {
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
		ch <- ReadingResult{
			Source: w.ctx.SensorID,
			Value:  value,
		}
	}
}

// WriteWithError allows to Write reading results and logged potential error at the same time.
func (w *MetricWriter) WriteWithError(value interface{}, err error) {
	if err != nil {
		w.ctx.Error(err)
		return
	}

	w.Write(value)
}
