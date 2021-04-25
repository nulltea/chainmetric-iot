package engine

import (
	"context"

	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/model"
)

type ReceiverFunc func(model.SensorsReadingResults)

type Request struct {
	Context context.Context
	Metrics []models.Metric
	Handler ReceiverFunc
}

