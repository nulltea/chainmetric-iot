package events

const (
	// DeviceLocationChanged identifies event for changes in location of the models.Device.
	DeviceLocationChanged = "device.location.changed"

	// MetricReadingsPostFailed identifies event for failure of models.MetricReadings post.
	MetricReadingsPostFailed = "readings.post.failed"

	// RequirementsSubmitted identifies event for submitting or changing models.Requirements request.
	RequirementsSubmitted = "requirements.submitted"

	// SensorsRegisterChanged identifies event for changes in sensor.SensorsRegister.
	SensorsRegisterChanged = "sensors.register.changed"

	// DeviceLoggedOnNetwork identifies event for the models.Device to became logged on network.
	DeviceLoggedOnNetwork = "device.logged"

	// DeviceRemovedFromNetwork identifies event for the removal of the models.Device from network.
	DeviceRemovedFromNetwork = "device.removed"
)
