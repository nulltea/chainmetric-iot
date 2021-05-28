package events

const (
	// DeviceLocationChanged identifies event for changes in location of the models.Device.
	DeviceLocationChanged = "device.location.changed"

	// MetricReadingsPostFailed identifies event for failure of models.MetricReadings post.
	MetricReadingsPostFailed = "readings.post.failed"

	// RequirementsSubmitted identifies event for submitting or changing models.Requirements request.
	RequirementsSubmitted = "requirements.submitted"
)
