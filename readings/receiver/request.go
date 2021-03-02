package receiver

import (
	"sensorsys/model"
)

type ReceiverFunc func(model.MetricReadings)

type Request struct {
	Metrics []model.Metric
	Handler ReceiverFunc
}

