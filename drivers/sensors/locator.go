package sensors

import (
	"github.com/spf13/viper"

	"github.com/timoth-y/chainmetric-sensorsys/core/sensor"
)

var i2cSensorsLocatorMap = map[uint16][]sensor.Factory {
	0x1D: { sensor.I2CFactory(NewAccelerometerLSM303, LSM303C_A_ADDRESS) },
	0x1E: { sensor.I2CFactory(NewMagnetometerLSM303, LSM303C_M_ADDRESS) },
	0x40: { sensor.I2CFactory(NewHDC1080, HDC1080_ADDRESS) },
	0x48: { sensor.I2CFactory(NewADCHall, ADC_HALL_ADDRESS) },
	0x49: { sensor.I2CFactory(NewADCMicrophone, ADC_MICROPHONE_ADDRESS) },
	0x4A: {
		sensor.I2CFactory(NewMAX44009, MAX44009_ADDRESS),
		sensor.I2CFactory(NewADCMQ9, ADC_MQ9_ADDRESS),
	},
	0x4B: {
		sensor.I2CFactory(NewMAX44009, MAX44009_ALT_ADDRESS),
		sensor.I2CFactory(NewADCPiezo, ADC_PIEZO_ADDRESS),
	},
	0x53: { sensor.I2CFactory(NewADXL345, ADXL345_ADDRESS) },
	0x57: { sensor.I2CFactory(NewMAX30102, MAX30102_ADDRESS) },
	0x5A: { sensor.I2CFactory(NewCCS811, CCS811_ADDRESS) },
	0x60: { sensor.I2CFactory(NewSI1145, SI1145_ADDRESS) },
	0x76: { sensor.I2CFactory(NewBMXX80, BMP280_ADDRESS) },
	0x88: { sensor.I2CFactory(NewI2CSensorMock, MOCK_ADDRESS) },
}

// LocateI2CSensor locates I2C-based sensor.Sensor and provides its sensor.Factory.
func LocateI2CSensor(addr uint16, bus int) (sensor.Factory, bool) {
	if factories, ok := i2cSensorsLocatorMap[addr]; ok {
		for i, f := range factories {
			if f.Build(bus).Verify() {
				return factories[i], true
			}
		}
	}

	return nil, false
}

// I2CAddressesRange determines diapason of I2C addresses to detect from.
func I2CAddressesRange() []uint16 {
	var addresses []uint16

	for addr := range i2cSensorsLocatorMap {
		if addr == MOCK_ADDRESS && !viper.GetBool("mocks.debug_env") {
			continue
		}

		addresses = append(addresses, addr)
	}

	return addresses
}
