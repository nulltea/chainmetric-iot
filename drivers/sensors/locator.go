package sensors

import (
	"sort"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/sensor"
)

var i2cSensorsLocatorMap = map[uint16]sensor.Factory{
	ADXL345_ADDRESS:        sensor.I2CFactory(ADXL345_ADDRESS, NewADXL345),
	BMP280_ADDRESS:         sensor.I2CFactory(BMP280_ADDRESS, NewBMXX80),
	CCS811_ADDRESS:         sensor.I2CFactory(CCS811_ADDRESS, NewCCS811),
	HDC1080_ADDRESS:        sensor.I2CFactory(HDC1080_ADDRESS, NewHDC1080),
	MAX30102_ADDRESS:       sensor.I2CFactory(MAX30102_ADDRESS, NewMAX30102),
	MAX44009_ADDRESS:       sensor.I2CFactory(MAX44009_ADDRESS, NewMAX44009),
	SI1145_ADDRESS:         sensor.I2CFactory(SI1145_ADDRESS, NewSI1145),
	LSM303C_A_ADDRESS:      sensor.I2CFactory(LSM303C_A_ADDRESS, NewAccelerometerLSM303),
	LSM303C_M_ADDRESS:      sensor.I2CFactory(LSM303C_M_ADDRESS, NewMagnetometerLSM303),
	ADC_HALL_ADDRESS:       sensor.I2CFactory(ADC_HALL_ADDRESS, NewADCHall),
	ADC_PIEZO_ADDRESS:      sensor.I2CFactory(ADC_PIEZO_ADDRESS, NewADCPiezo),
	ADC_MICROPHONE_ADDRESS: sensor.I2CFactory(ADC_MICROPHONE_ADDRESS, NewADCMicrophone),
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

	return uint16(addresses[0]), uint16(len(addresses) - 1)
}
