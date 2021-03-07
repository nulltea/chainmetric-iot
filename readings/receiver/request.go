package receiver

import (
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/model"
)

type ReceiverFunc func(model.MetricReadings)

type Request struct {
	Metrics []models.Metric
	Handler ReceiverFunc
}

