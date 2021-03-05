package config

type SensorsConfig struct {
	ADXL345  I2CSensorConfig        `yaml:"ADXL345"`
	BMP280   I2CSensorConfig        `yaml:"BMP280"`
	CCS811   I2CSensorConfig        `yaml:"BMP280"`
	DHT22    DigitalPinSensorConfig `yaml:"DHT22"`
	MAX30102 I2CSensorConfig        `yaml:"MAX30102"`
	MAX44009 I2CSensorConfig        `yaml:"MAX44009"`
	SI1145   I2CSensorConfig        `yaml:"SI1145"`
}

type SensorCommonConfig struct {
	Timeout int `yaml:"timeout"`
}

type I2CSensorConfig struct {
	SensorCommonConfig
	Address uint8 `yaml:"address"`
	Bus     int   `yaml:"bus"`
}

type SPISensorConfig struct {
	SensorCommonConfig
}

type DigitalPinSensorConfig struct {
	SensorCommonConfig
	Pin int `yaml:"pin"`
}

type AnalogPinSensorConfig struct {
	SensorCommonConfig
	Pin int `yaml:"pin"`
}
