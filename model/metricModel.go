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
	AirC02Concentration Metric = "air_C02_cc"
	AirTVOCsConcentration Metric = "air_TVOC_cc"
)
