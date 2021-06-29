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
	cancel  context.CancelFunc
}

// Cancel calls assigned cancel func to cancel request receiver routine.
func (sr *SensorsReadingRequest) Cancel() {
	if sr.IsProcessed() {
		sr.cancel()
	}
}

// SetCancel sets `cancel` func for canceling request receiver routine.
func (sr *SensorsReadingRequest) SetCancel(cancel context.CancelFunc) {
	sr.cancel = cancel
}

// IsProcessed determines whether the request is being already processed.
func (sr *SensorsReadingRequest) IsProcessed() bool {
	return sr.cancel != nil
}

