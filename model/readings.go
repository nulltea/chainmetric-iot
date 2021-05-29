package model

import (
	"context"
	"time"

	"github.com/timoth-y/chainmetric-core/models"
)

type SensorReading struct {
	Source string
	Value float64
}

type SensorsReadingResults map[models.Metric] float64

type SensorReadingsPipe map[models.Metric] chan SensorReading

 // SensorsReadingRequest defines structure of the SensorReading request based on models.Requirements record.
type SensorsReadingRequest struct {
	ID      string
	AssetID string
	Period  time.Duration
	Metrics models.Metrics
	Cancel  context.CancelFunc
}
