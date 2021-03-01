package readings

import (
	"sensorsys/model"
	"sensorsys/readings/receiver"
)

type Request struct {
	Metrics []model.Metric
	Context *receiver.Context
}

