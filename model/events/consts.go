package events

const (
	// DeviceLocationChanged identifies event for changes in location of the models.Device.
	DeviceLocationChanged = "device.location.changed"

	// MetricReadingsPostFailed identifies event for failure of models.MetricReadings post.
	MetricReadingsPostFailed = "readings.post.failed"

	// RequirementsChanged identifies event for submitting or changing models.Requirements request.
	RequirementsChanged = "requirements.changed"

	// AssetsChanged identifies event for submitting or changing models.Asset's.
	AssetsChanged = "assets.changed"

	// SensorsRegisterChanged identifies event for changes in sensor.SensorsRegister.
	SensorsRegisterChanged = "sensors.register.changed"

	// DeviceLoggedOnNetwork identifies event for the models.Device to became logged on network.
	DeviceLoggedOnNetwork = "device.logged"

	// DeviceRemovedFromNetwork identifies event for the removal of the models.Device from network.
	DeviceRemovedFromNetwork = "device.removed"

	// CacheChanged identifies event changes parameters cache.
	CacheChanged = "cache.changed"
)
