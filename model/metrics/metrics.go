package metrics

import "github.com/timoth-y/iot-blockchain-contracts/models"

var (
	Temperature           models.Metric = "temperature"
	Humidity              models.Metric = "humidity"
	Luminosity            models.Metric = "luminosity"
	Magnetism             models.Metric = "magnetism"
	Pressure              models.Metric = "pressure"
	Altitude              models.Metric = "altitude"
	UVLight               models.Metric = "uv_light"
	VisibleLight          models.Metric = "visible_light"
	IRLight               models.Metric = "ir_light"
	Proximity             models.Metric = "proximity"
	AirCO2Concentration   models.Metric = "air_CO2_concentration"
	AirTVOCsConcentration models.Metric = "air_TVOC_concentration"
	AccelerationInG       models.Metric = "acceleration_axes_G"
	AccelerationInMS2     models.Metric = "acceleration_axes_ms/2"
	HeartRate             models.Metric = "heart_rate"
	BloodOxidation        models.Metric = "blood_oxidation"
)
