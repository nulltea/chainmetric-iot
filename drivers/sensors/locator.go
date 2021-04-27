package sensors

import (
	"sort"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/periphery"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
)

var (
	ADXL345_KEY        = periphery.I2CDetectKey{ADXL345_ADDRESS, ADXL345_DEVICE_ID_REGISTER, ADXL345_DEVICE_ID}
	BMP280_KEY         = periphery.I2CDetectKey{BMP280_ADDRESS, BMP280_DEVICE_ID_REGISTER, BMP280_DEVICE_ID}
	CCS811_KEY         = periphery.I2CDetectKey{CCS811_ADDRESS, CCS811_DEVICE_ID_REGISTER, CCS811_DEVICE_ID}
	HDC1080_KEY        = periphery.I2CDetectKey{HDC1080_ADDRESS, HDC1080_DEVICE_ID_REGISTER, HDC1080_DEVICE_ID}
	MAX30102_KEY       = periphery.I2CDetectKey{MAX30102_ADDRESS, MAX30102_DEVICE_ID_REGISTER, MAX30102_DEVICE_ID}
	MAX44009_KEY       = periphery.I2CDetectKey{MAX44009_ADDRESS, MAX44009_DEVICE_ID_REGISTER, MAX44009_DEVICE_ID}
	SI1145_KEY         = periphery.I2CDetectKey{SI1145_ADDRESS, SI1145_DEVICE_ID_REGISTER, SI1145_DEVICE_ID}
	LSM303C_A_KEY      = periphery.I2CDetectKey{LSM303C_A_ADDRESS, LSM303C_A_DEVICE_ID_REGISTER, LSM303C_A_DEVICE_ID}
	LSM303C_M_KEY      = periphery.I2CDetectKey{LSM303C_M_ADDRESS, LSM303C_M_DEVICE_ID_REGISTER, LSM303C_M_DEVICE_ID}
	ADC_HALL_KEY       = periphery.I2CDetectKey{ADC_HALL_ADDRESS, ADS1115_DEVICE_ID_REGISTER, ADS1115_DEVICE_ID}
	ADC_MICROPHONE_KEY = periphery.I2CDetectKey{ADC_MICROPHONE_ADDRESS, ADS1115_DEVICE_ID_REGISTER, ADS1115_DEVICE_ID}
	ADC_PIEZO_KEY      = periphery.I2CDetectKey{ADC_PIEZO_ADDRESS, ADS1115_DEVICE_ID_REGISTER, ADS1115_DEVICE_ID}
	ADC_FLAME_KEY      = periphery.I2CDetectKey{ADC_FLAME_ADDRESS, ADS1115_DEVICE_ID_REGISTER, ADS1115_DEVICE_ID}
	ADC_MQ9_KEY        = periphery.I2CDetectKey{ADC_MQ9_ADDRESS, ADS1115_DEVICE_ID_REGISTER, ADS1115_DEVICE_ID}
	MOCK_KEY           = periphery.I2CDetectKey{MOCK_ADDRESS, MOCK_DEVICE_ID_REGISTER, MOCK_DEVICE_ID}
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
