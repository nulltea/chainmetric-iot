package sensors

const (
	ADXL345_ADDRESS        = 0x53
	BMP280_ADDRESS         = 0x76
	CCS811_ADDRESS         = 0x5A
	HDC1080_ADDRESS        = 0x40
	MAX30102_ADDRESS       = 0x57
	MAX44009_ADDRESS       = 0x4A
	MAX44009_ALT_ADDRESS   = 0x4B
	SI1145_ADDRESS         = 0x60
	LSM303C_A_ADDRESS      = 0x1D
	LSM303C_M_ADDRESS      = 0x1E
	ADC_HALL_ADDRESS       = 0x48
	ADC_MICROPHONE_ADDRESS = 0x49
	ADC_MQ9_ADDRESS        = 0x4A
	ADC_FLAME_ADDRESS      = 0x4B
	ADC_PIEZO_ADDRESS      = 0x4E
	MOCK_ADDRESS           = 0x88
)

// ADCMicrophone sensor constants
const (
	ADC_MICROPHONE_BIAS          = 2500
	ADC_MICROPHONE_REGRESSION_C1 = 0.001276
	ADC_MICROPHONE_REGRESSION_C2 = 47.56
)

// ADCHall sensor constants
const (
	ADC_HALL_BIAS        = 400
	ADC_HALL_SENSITIVITY = 1.9 // mV / Gauss
)

// ADCFlame sensor constants
const (
	ADC_FLAME_BIAS = 2
)

// ADCMQ9 sensor constants
const (
	ADC_MQ9_BIAS        = -50
	ADC_MQ9_RESISTANCE  = 5
	ADC_MQ9_SENSITIVITY = 9.9
)

// ADCPiezo sensor constants
const (
	ADC_PIEZO_BIAS = 0
)

// ADXL345 accelerometer sensor constants
const (
	// Registers
	ADXL345_DEVICE_ID_REGISTER = 0x00
	ADXL345_DATA_FORMAT        = 0x31
	ADXL345_BW_RATE            = 0x2C
	ADXL345_POWER_CTL          = 0x2D
	ADXL345_MEASURE            = 0x08

	// Constants
	ADXL345_DEVICE_ID = 0xE5

	// Device bandwidth and output data rates
	ADXL345_Rate1600HZ = 0x0F
	ADXL345_Rate800HZ  = 0x0E
	ADXL345_Rate400HZ  = 0x0D
	ADXL345_Rate200HZ  = 0x0C
	ADXL345_Rate100HZ  = 0x0B
	ADXL345_Rate50HZ   = 0x0A
	ADXL345_Rate25HZ   = 0x09

	// Measurement Range
	ADXL345_RANGE2G  = 0x00
	ADXL345_RANGE4G  = 0x01
	ADXL345_RANGE8G  = 0x02
	ADXL345_RANGE16G = 0x03

	// Axes Data
	ADXL345_DATAX0 = 0x32
	ADXL345_DATAX1 = 0x33
	ADXL345_DATAY0 = 0x34
	ADXL345_DATAY1 = 0x35
	ADXL345_DATAZ0 = 0x36
	ADXL345_DATAZ1 = 0x37
)

// BMP280 barometer sensor constants
const (
	BMP280_DEVICE_ID_REGISTER = 0xD0
	BMP280_DEVICE_ID = 0x60
)

// CCS811 air quality sensor constants
const (
	// Registers
	CCS811_STATUS             = 0x00
	CCS811_MEAS_MODE          = 0x01
	CCS811_ALG_RESULT_DATA    = 0x02
	CCS811_RAW_DATA           = 0x03
	CCS811_ENV_DATA           = 0x05
	CCS811_NTC                = 0x06
	CCS811_THRESHOLDS         = 0x10
	CCS811_BASELINE           = 0x11
	CCS811_DEVICE_ID_REGISTER = 0x20
	CCS811_HW_VERSION         = 0x21
	CCS811_FW_BOOT_VERSION    = 0x23
	CCS811_FW_APP_VERSION     = 0x24
	CCS811_ERROR_ID           = 0xE0
	CCS811_SW_RESET           = 0xFF

	// Constants
	CCS811_DEVICE_ID    = 0x81
	CCS811_REF_RESISTOR = 100000

	// Bootloader Registers
	CCS811_BOOTLOADER_APP_ERASE  = 0xF1
	CCS811_BOOTLOADER_APP_DATA   = 0xF2
	CCS811_BOOTLOADER_APP_VERIFY = 0xF3
	CCS811_BOOTLOADER_APP_START  = 0xF4

	// Drive mode
	CCS811_DRIVE_MODE_IDLE  = 0x00
	CCS811_DRIVE_MODE_1SEC  = 0x01
	CCS811_DRIVE_MODE_10SEC = 0x02
	CCS811_DRIVE_MODE_60SEC = 0x03
	CCS811_DRIVE_MODE_250MS = 0x04

	// STATUS - Bitwise
	CCS811_ERROR_BIT      = 0x01
	CCS811_DATA_READY_BIT = 0x08
	CCS811_APP_VALID_BIT  = 0x10
	CCS811_FW_MODE_BIT    = 0x80

	// ERROR - Bitwise
	CCS811_WRITE_REG_INVALID = 0x01
	CCS811_READ_REG_INVALID  = 0x02
	CCS811_MEASMODE_INVALID  = 0x04
	CCS811_MAX_RESISTANCE    = 0x08
	CCS811_HEATER_FAULT      = 0x10
	CCS811_HEATER_SUPPLY     = 0x20

	// Time
	CCS811_APP_START_TIME    = 100
	CCS811_RESET_TIME    = 100
	CCS811_RETRY_TIME = 250
)

// HDC1080 temperature and humidity sensor constants
const (
	// Registers
	HDC1080_TEMPERATURE_REGISTER =          0x00
	HDC1080_HUMIDITY_REGISTER =             0x01
	HDC1080_CONFIGURATION_REGISTER =        0x02
	HDC1080_DEVICE_ID_REGISTER =            0xFF

	// Device ID
	HDC1080_DEVICE_ID = 0x10

	// Configuration Register Bits
	HDC1080_CONFIG_RESET_BIT =                0x8000
	HDC1080_CONFIG_ACQUISITION_MODE =         0x1000
	HDC1080_CONFIG_TEMPERATURE_RESOLUTION =   0x0400
	HDC1080_CONFIG_HUMIDITY_RESOLUTION_HBIT = 0x0200
	HDC1080_CONFIG_HUMIDITY_RESOLUTION_LBIT = 0x0100
)

// LSM303Accelerometer and LSM303Magnetometer sensors constants
const (
	// Registers
	LSM303C_A_DEVICE_ID_REGISTER = 0x0F
	LSM303C_M_DEVICE_ID_REGISTER = 0x0F

	// Constants
	LSM303C_A_DEVICE_ID = 0x41
	LSM303C_M_DEVICE_ID = 0x3D
)

// MAX30102 pulse-oximeter sensor constants
const(
	MAX30102_DEVICE_ID_REGISTER = 0xFF
	MAX30102_DEVICE_ID = 0x15
)

// MAX44009 luminosity sensor constants
const(
	// Commands
	MAX44009_APP_START = 0x03

	// Registers
	MAX44009_DEVICE_ID_REGISTER = 0x0F

	// Constants
	MAX44009_DEVICE_ID = 0x3F
)

// SI1145 ambient light sensor constants
const(
	// Commands
	SI1145_PARAM_QUERY = 0x80
	SI1145_PARAM_SET   = 0xA0
	SI1145_NOP         = 0x0
	SI1145_RESET       = 0x01
	SI1145_BUSADDR     = 0x02
	SI1145_PS_FORCE    = 0x05
	SI1145_ALS_FORCE   = 0x06
	SI1145_PSALS_FORCE = 0x07
	SI1145_PS_PAUSE    = 0x09
	SI1145_ALS_PAUSE   = 0x0A
	SI1145_PSALS_PAUSE = 0xB
	SI1145_PS_AUTO     = 0x0D
	SI1145_ALS_AUTO    = 0x0E
	SI1145_PSALS_AUTO  = 0x0F
	SI1145_GET_CAL     = 0x12

	SI1145_DEVICE_ID_REGISTER = 0x02
	SI1145_DEVICE_ID = 0x08

	// Parameters
	SI1145_PARAM_I2CADDR         = 0x00
	SI1145_PARAM_CHLIST          = 0x01
	SI1145_PARAM_CHLIST_ENUV     = 0x80
	SI1145_PARAM_CHLIST_ENAUX    = 0x40
	SI1145_PARAM_CHLIST_ENALSIR  = 0x20
	SI1145_PARAM_CHLIST_ENALSVIS = 0x10
	SI1145_PARAM_CHLIST_ENPS1    = 0x01
	SI1145_PARAM_CHLIST_ENPS2    = 0x02
	SI1145_PARAM_CHLIST_ENPS3    = 0x04

	SI1145_PARAM_PSLED12SEL         = 0x02
	SI1145_PARAM_PSLED12SEL_PS2NONE = 0x00
	SI1145_PARAM_PSLED12SEL_PS2LED1 = 0x10
	SI1145_PARAM_PSLED12SEL_PS2LED2 = 0x20
	SI1145_PARAM_PSLED12SEL_PS2LED3 = 0x40
	SI1145_PARAM_PSLED12SEL_PS1NONE = 0x00
	SI1145_PARAM_PSLED12SEL_PS1LED1 = 0x01
	SI1145_PARAM_PSLED12SEL_PS1LED2 = 0x02
	SI1145_PARAM_PSLED12SEL_PS1LED3 = 0x04

	SI1145_PARAM_PSLED3SEL = 0x03
	SI1145_PARAM_PSENCODE  = 0x05
	SI1145_PARAM_ALSENCODE = 0x06

	SI1145_PARAM_PS1ADCMUX        = 0x07
	SI1145_PARAM_PS2ADCMUX        = 0x08
	SI1145_PARAM_PS3ADCMUX        = 0x09
	SI1145_PARAM_PSADCOUNTER      = 0x0A
	SI1145_PARAM_PSADCGAIN        = 0x0B
	SI1145_PARAM_PSADCMISC        = 0x0C
	SI1145_PARAM_PSADCMISC_RANGE  = 0x20
	SI1145_PARAM_PSADCMISC_PSMODE = 0x04

	SI1145_PARAM_ALSIRADCMUX = 0x0E
	SI1145_PARAM_AUXADCMUX   = 0x0F

	SI1145_PARAM_ALSVISADCOUNTER        = 0x10
	SI1145_PARAM_ALSVISADCGAIN          = 0x11
	SI1145_PARAM_ALSVISADCMISC          = 0x12
	SI1145_PARAM_ALSVISADCMISC_VISRANGE = 0x20

	SI1145_PARAM_ALSIRADCOUNTER     = 0x1D
	SI1145_PARAM_ALSIRADCGAIN       = 0x1E
	SI1145_PARAM_ALSIRADCMISC       = 0x1F
	SI1145_PARAM_ALSIRADCMISC_RANGE = 0x20

	SI1145_PARAM_ADCCOUNTER_511CLK = 0x70

	SI1145_PARAM_ADCMUX_SMALLIR = 0x00
	SI1145_PARAM_ADCMUX_LARGEIR = 0x03

	// REGISTERS
	SI1145_REG_PARTID = 0x00
	SI1145_REG_REVID  = 0x01
	SI1145_REG_SEQID  = 0x02

	SI1145_REG_INTCFG         = 0x03
	SI1145_REG_INTCFG_INTOE   = 0x01
	SI1145_REG_INTCFG_INTMODE = 0x02

	SI1145_REG_IRQEN                = 0x04
	SI1145_REG_IRQEN_ALSEVERYSAMPLE = 0x01
	SI1145_REG_IRQEN_PS1EVERYSAMPLE = 0x04
	SI1145_REG_IRQEN_PS2EVERYSAMPLE = 0x08
	SI1145_REG_IRQEN_PS3EVERYSAMPLE = 0x10

	SI1145_REG_IRQMODE1 = 0x05
	SI1145_REG_IRQMODE2 = 0x06

	SI1145_REG_HWKEY       = 0x07
	SI1145_REG_MEASRATE0   = 0x08
	SI1145_REG_MEASRATE1   = 0x09
	SI1145_REG_PSRATE      = 0x0A
	SI1145_REG_PSLED21     = 0x0F
	SI1145_REG_PSLED3      = 0x10
	SI1145_REG_UCOEFF0     = 0x13
	SI1145_REG_UCOEFF1     = 0x14
	SI1145_REG_UCOEFF2     = 0x15
	SI1145_REG_UCOEFF3     = 0x16
	SI1145_REG_PARAMWR     = 0x17
	SI1145_REG_COMMAND     = 0x18
	SI1145_REG_RESPONSE    = 0x20
	SI1145_REG_IRQSTAT     = 0x21
	SI1145_REG_IRQSTAT_ALS = 0x01

	SI1145_REG_ALSVISDATA0 = 0x22
	SI1145_REG_ALSVISDATA1 = 0x23
	SI1145_REG_ALSIRDATA0  = 0x24
	SI1145_REG_ALSIRDATA1  = 0x25
	SI1145_REG_PS1DATA0    = 0x26
	SI1145_REG_PS1DATA1    = 0x27
	SI1145_REG_PS2DATA0    = 0x28
	SI1145_REG_PS2DATA1    = 0x29
	SI1145_REG_PS3DATA0    = 0x2A
	SI1145_REG_PS3DATA1    = 0x2B
	SI1145_REG_UVINDEX0    = 0x2C
	SI1145_REG_UVINDEX1    = 0x2D
	SI1145_REG_PARAMRD     = 0x2E
	SI1145_REG_CHIPSTAT    = 0x30
)

const (
	MOCK_DEVICE_ID_REGISTER = 0x0F
	MOCK_DEVICE_ID          = 0x69
)
