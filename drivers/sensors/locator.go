package sensors

const (
	ADXL345_ADDRESS        = 0x53
	BMP280_ADDRESS         = 0x76
	CCS811_ADDRESS         = 0x5A
	HDC1080_ADDRESS        = 0x40
	MAX30102_ADDRESS       = 0x57
	MAX44009_ADDRESS       = 0x4B
	SI1145_ADDRESS         = 0x60
	LSM303C_A_ADDRESS      = 0x1D
	LSM303C_M_ADDRESS      = 0x1E
	ADC_HALL_ADDRESS       = 0x48
	ADC_PIEZO_ADDRESS      = 0x49
	ADC_MICROPHONE_ADDRESS = 0x4A
)

var I2CSensorsMap = map[uint16]func(bus int) Sensor{
	ADXL345_ADDRESS: func(bus int) Sensor {
		return NewADXL345(ADXL345_ADDRESS, bus)
	},
	BMP280_ADDRESS: func(bus int) Sensor {
		return NewBMXX80(BMP280_ADDRESS, bus)
	},
	CCS811_ADDRESS: func(bus int) Sensor {
		return NewCCS811(CCS811_ADDRESS, bus)
	},
	HDC1080_ADDRESS: func(bus int) Sensor {
		return NewHDC1080(HDC1080_ADDRESS, bus)
	},
	MAX30102_ADDRESS: func(bus int) Sensor {
		return NewMAX30102(MAX44009_ADDRESS, bus)
	},
	MAX44009_ADDRESS: func(bus int) Sensor {
		return NewMAX44009(MAX44009_ADDRESS, bus)
	},
	SI1145_ADDRESS: func(bus int) Sensor {
		return NewSI1145(SI1145_ADDRESS, bus)
	},
	LSM303C_A_ADDRESS: func(bus int) Sensor {
		return NewAccelerometerLSM303(LSM303C_A_ADDRESS, bus)
	},
	LSM303C_M_ADDRESS: func(bus int) Sensor {
		return NewMagnetometerLSM303(LSM303C_M_ADDRESS, bus)
	},
	ADC_HALL_ADDRESS: func(bus int) Sensor {
		return NewADCHall(ADC_HALL_ADDRESS, bus)
	},
	ADC_PIEZO_ADDRESS: func(bus int) Sensor {
		return NewADCPiezo(ADC_PIEZO_ADDRESS, bus)
	},
	ADC_MICROPHONE_ADDRESS: func(bus int) Sensor {
		return NewADCMicrophone(ADC_MICROPHONE_ADDRESS, bus)
	},
}
