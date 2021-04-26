package sensors

import (
	"sort"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
)

var i2cSensorsLocatorMap = map[uint16]sensor.Factory{
	ADXL345_ADDRESS:        sensor.I2CFactory(NewADXL345, ADXL345_ADDRESS),
	BMP280_ADDRESS:         sensor.I2CFactory(NewBMXX80, BMP280_ADDRESS),
	CCS811_ADDRESS:         sensor.I2CFactory(NewCCS811, CCS811_ADDRESS),
	HDC1080_ADDRESS:        sensor.I2CFactory(NewHDC1080, HDC1080_ADDRESS),
	MAX30102_ADDRESS:       sensor.I2CFactory(NewMAX30102, MAX30102_ADDRESS),
	MAX44009_ADDRESS:       sensor.I2CFactory(NewMAX44009, MAX44009_ADDRESS),
	MAX44009_ALT_ADDRESS:   sensor.I2CFactory(NewMAX44009, MAX44009_ADDRESS),
	SI1145_ADDRESS:         sensor.I2CFactory(NewSI1145, SI1145_ADDRESS),
	LSM303C_A_ADDRESS:      sensor.I2CFactory(NewAccelerometerLSM303, LSM303C_A_ADDRESS),
	LSM303C_M_ADDRESS:      sensor.I2CFactory(NewMagnetometerLSM303, LSM303C_M_ADDRESS),
	ADC_HALL_ADDRESS:       sensor.I2CFactory(NewADCHall, ADC_HALL_ADDRESS),
	ADC_MICROPHONE_ADDRESS: sensor.I2CFactory(NewADCMicrophone, ADC_MICROPHONE_ADDRESS),
	ADC_PIEZO_ADDRESS:      sensor.I2CFactory(NewADCPiezo, ADC_PIEZO_ADDRESS),
	ADC_FLAME_ADDRESS:      sensor.I2CFactory(NewADCFlame, ADC_FLAME_ADDRESS),
	ADC_MQ9_ADDRESS:        sensor.I2CFactory(NewADCMQ9, ADC_MQ9_ADDRESS),
	MOCK_ADDRESS:           sensor.I2CFactory(NewI2CSensorMock, MOCK_ADDRESS),
}

// LocateI2CSensor locates I2C-based sensor.Sensor and provides its sensor.Factory.
func LocateI2CSensor(addr uint16) (f sensor.Factory, ok bool) {
	f, ok = i2cSensorsLocatorMap[addr]
	return
}

// I2CAddressesRange determines diapason of I2C addresses to detect from.
func I2CAddressesRange() (from uint16, to uint16) {
	var addresses []int

	for addr := range i2cSensorsLocatorMap {
		addresses = append(addresses, int(addr))
	}

	sort.Ints(addresses)

	return uint16(addresses[0]), uint16(addresses[len(addresses) - 1])
}
