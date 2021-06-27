package events

import (
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-iot/core/dev/sensor"
	"github.com/timoth-y/chainmetric-iot/model"
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

// AssetsChangedPayload defines payload for AssetsChanged event.
type AssetsChangedPayload struct {
	Assigned []string
	Removed  []string
}

// RequirementsChangedPayload defines payload for RequirementsChanged event.
type RequirementsChangedPayload struct {
	Requests []model.SensorsReadingRequest
}

// SensorsRegisterChangedPayload defines payload for SensorsRegisterChanged event.
type SensorsRegisterChangedPayload struct {
	Added   []sensor.Sensor
	Removed []string
}
