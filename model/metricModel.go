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

var (
	Temperature           Metric = "temperature"
	Humidity              Metric = "humidity"
	Luminosity            Metric = "luminosity"
	Magnetism             Metric = "magnetism"
	Pressure              Metric = "pressure"
	Altitude              Metric = "altitude"
	UVLight               Metric = "uv_light"
	VisibleLight          Metric = "visible_light"
	IRLight               Metric = "ir_light"
	AirCO2Concentration   Metric = "air_CO2_concentration"
	AirTVOCsConcentration Metric = "air_TVOC_concentration"
	Acceleration          Metric = "acceleration"
	HeartRate             Metric = "heart_rate"
	BloodOxidation        Metric = "blood_oxidation"
)
