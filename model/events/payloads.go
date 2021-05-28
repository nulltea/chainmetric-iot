package events

import "github.com/timoth-y/chainmetric-core/models"

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
