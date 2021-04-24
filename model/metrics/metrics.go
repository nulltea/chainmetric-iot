package metrics

import "github.com/timoth-y/iot-blockchain-contracts/models"

var (
	Temperature               models.Metric = "temp"
	Humidity                  models.Metric = "hdt"
	Luminosity                models.Metric = "lux"
	Magnetism                 models.Metric = "mgn"
	Pressure                  models.Metric = "bar"
	Altitude                  models.Metric = "alt"
	UVLight                   models.Metric = "uv"
	VisibleLight              models.Metric = "vis"
	IRLight                   models.Metric = "ir"
	Proximity                 models.Metric = "prx"
	AirCO2Concentration       models.Metric = "co2"
	AirTVOCsConcentration     models.Metric = "tvoc"
	AirPetroleumConcentration models.Metric = "lpg"
	Acceleration              models.Metric = "axg"
	HeartRate                 models.Metric = "hrt"
	BloodOxidation            models.Metric = "blo"
	Vibration                 models.Metric = "vbr"
	NoiseLevel                models.Metric = "nse"
	Flame                     models.Metric = "flm"
)
