package engine

import (
	"context"

	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-sensorsys/model"
)

type ReceiverFunc func(model.SensorsReadingResults)

type Request struct {
	Context context.Context
	Metrics []models.Metric
	Handler ReceiverFunc
}

