package model

type Metric string

// temperature - DHT11, 18B20
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
)