package model

type Metric string

// temperature - DHT22,
// humidity - DHT11
// vibration - KY002, KY031
// sound
// magnetic - hall effect sensors
// flame - KY026
// luminosity - TEMT6000
// gas
// accelerometer - ADXL345

var(
	Temperature Metric = "temperature"
	Humidity Metric = "humidity"
	Luminosity Metric = "luminosity"
	Magnetism Metric = "magnetism"
	Pressure Metric = "pressure"
	Altitude Metric = "altitude"
	UVLight Metric = "uv_light"
	VisibleLight Metric = "visible_light"
	IRLight Metric = "ir_light"
	AirC02Concentration Metric = "air_C02_concentration"
	AirTVOCsConcentration Metric = "air_TVOC_concentration"
)
