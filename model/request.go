package model

import (
	"context"
	"time"

	"github.com/timoth-y/chainmetric-core/models"
)

// SensorsReadingRequest defines structure of the ReadingResult request based on models.Requirements record.
type SensorsReadingRequest struct {
	ID      string
	AssetID string
	Period  time.Duration
	Metrics models.Metrics
	Cancel  context.CancelFunc
}
