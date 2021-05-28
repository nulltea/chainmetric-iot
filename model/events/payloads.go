package events

import (
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-sensorsys/model"
)

// DeviceLocationChangedPayload defines payload for DeviceLocationChanged event.
type DeviceLocationChangedPayload struct {
	Old models.Location
	New models.Location
}

// MetricReadingsPostFailedPayload defines payload for MetricReadingsPostFailed event.
type MetricReadingsPostFailedPayload struct {
	models.MetricReadings
	Error error
}

// RequirementsSubmittedPayload defines payload for RequirementsSubmitted event.
type RequirementsSubmittedPayload struct {
	Requests []model.SensorsReadingRequest
}
