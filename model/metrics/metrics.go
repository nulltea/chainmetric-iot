package metrics

import "sensorsys/model"

var (
	Temperature           model.Metric = "temperature"
	Humidity              model.Metric = "humidity"
	Luminosity            model.Metric = "luminosity"
	Magnetism             model.Metric = "magnetism"
	Pressure              model.Metric = "pressure"
	Altitude              model.Metric = "altitude"
	UVLight               model.Metric = "uv_light"
	VisibleLight          model.Metric = "visible_light"
	IRLight               model.Metric = "ir_light"
	Proximity             model.Metric = "proximity"
	AirCO2Concentration   model.Metric = "air_CO2_concentration"
	AirTVOCsConcentration model.Metric = "air_TVOC_concentration"
	AccelerationInG       model.Metric = "acceleration_axes_G"
	AccelerationInMS2     model.Metric = "acceleration_axes_ms/2"
	HeartRate             model.Metric = "heart_rate"
	BloodOxidation        model.Metric = "blood_oxidation"
)
