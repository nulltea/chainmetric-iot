package sensors

const (
	BMP280_ADDRESS = 0x76
	CCS811_ADDRESS = 0x5A
	MAX30102_ADDRESS = 0x57
	MAX44009_ADDRESS = 0x4A
	SI1145_ADDRESS = 0x60
)

var (
	I2CSensorsMap = map[uint8]func(bus int) Sensor {
		BMP280_ADDRESS: func(bus int) Sensor {
			return NewBMP280(BMP280_ADDRESS, bus)
		},
		CCS811_ADDRESS: func(bus int) Sensor {
			return NewCCS811(CCS811_ADDRESS, bus)
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
	}
)
