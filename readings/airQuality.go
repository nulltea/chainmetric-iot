package readings

import (
	"context"
	"fmt"
	"time"

	"github.com/d2r2/go-i2c"

	"sensorsys/model"
)

const(
	//// Registers
	CCS811_STATUS = 0x00
	CCS811_MEAS_MODE = 0x01
	CCS811_ALG_RESULT_DATA = 0x02
	CCS811_RAW_DATA = 0x03
	CCS811_ENV_DATA = 0x05
	CCS811_NTC = 0x06
	CCS811_THRESHOLDS = 0x10
	CCS811_BASELINE = 0x11
	CCS811_HW_ID = 0x20
	CCS811_HW_VERSION = 0x21
	CCS811_FW_BOOT_VERSION = 0x23
	CCS811_FW_APP_VERSION = 0x24
	CCS811_ERROR_ID = 0xE0
	CCS811_SW_RESET = 0xFF
	//// Bootloader Registers
	CCS811_BOOTLOADER_APP_ERASE = 0xF1
	CCS811_BOOTLOADER_APP_DATA = 0xF2
	CCS811_BOOTLOADER_APP_VERIFY = 0xF3
	CCS811_BOOTLOADER_APP_START = 0xF4
	//// Drive mode
	CCS811_DRIVE_MODE_IDLE = 0x00
	CCS811_DRIVE_MODE_1SEC = 0x01
	CCS811_DRIVE_MODE_10SEC = 0x02
	CCS811_DRIVE_MODE_60SEC = 0x03
	CCS811_DRIVE_MODE_250MS = 0x04
	//// CONSTANTs
	CCS811_HW_ID_CODE	=	0x81
	CCS811_REF_RESISTOR	=	100000
	//// STATUS - Bitwise
	CCS811_ERROR_BIT      = 0x01
	CCS811_DATA_READY_BIT = 0x08
	CCS811_APP_VALID_BIT         = 0x10
	CCS811_FW_MODE_BIT           = 0x80
	//// ERROR - Bitwise
	CCS811_WRITE_REG_INVALID = 0x01
	CCS811_READ_REG_INVALID = 0x02
	CCS811_MEASMODE_INVALID = 0x04
	CCS811_MAX_RESISTANCE = 0x08
	CCS811_HEATER_FAULT = 0x10
	CCS811_HEATER_SUPPLY = 0x20
)

var (
	InterruptMode byte = 0
	InterruptThreshold byte = 0
	SamplingRate byte =  CCS811_DRIVE_MODE_10SEC
)

func (s *SensorsReader) SubscribeToAirQualityReadings(addr uint8, bus int) error {
	i2c, err := i2c.NewI2C(addr, bus); if err != nil {
		return err
	}
	s.cleanQueue = append(s.cleanQueue, i2c.Close)

	err = initCCS811(i2c); if err != nil {
		return err
	}

	s.subscribe(func(ctx context.Context) {
		eC02, eTVOC, err := readCCS811(i2c); if err != nil {
			fmt.Println(err)
		}
		s.readings <- model.MetricReadings{
			model.AirC02Concentration: eC02,
			model.AirTVOCsConcentration: eTVOC,
		}
		s.waitGroup.Done()
	})

	return nil
}

func initCCS811(i2c *i2c.I2C) error {
	if !verifyId(i2c) {
		return fmt.Errorf("not that sensor")
	}

	// setReset(i2c)
	// time.Sleep(100 * time.Millisecond)

	getStatus(i2c)

	_, err := i2c.WriteBytes([]byte{CCS811_BOOTLOADER_APP_START}); if err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)

	status, err := getStatus(i2c); if err != nil {
		return err
	}

	if status & CCS811_ERROR_BIT != 0 {
		return fmt.Errorf("CCS811 device has error")
	}

	if status & CCS811_FW_MODE_BIT == 0{
		return fmt.Errorf("CCS811 device is in FW mode")
	}

	// setConfig(i2c); if err != nil {
	// 	return err
	// }
	return nil
}

func readCCS811(i2c *i2c.I2C) (float32, float32, error){
	var eCO2, eTVOC uint16 = 0,  0

	for {
		ready, err := isDataReady(i2c); if err != nil {
			return 0, 0, err
		}
		if ready {
			data, _, err := i2c.ReadRegBytes(CCS811_ALG_RESULT_DATA, 4)
			if err != nil {
				return 0, 0, err
			}
			eCO2 = ( uint16(data[0]) << 8) | uint16(data[1])
			eTVOC = ( uint16(data[2]) << 8) | uint16(data[3])
			break
		}
		time.Sleep(1 * time.Second)
	}

	return float32(eCO2), float32(eTVOC), nil
}

func isDataReady(d *i2c.I2C) (bool, error) {
	sts, err := getStatus(d)
	if err != nil {
		return false, err
	}

	return (sts & CCS811_DATA_READY_BIT) != 0, nil
}

func getStatus(d *i2c.I2C) (byte, error) {
	data, _, err := d.ReadRegBytes(CCS811_STATUS, 1); if err != nil {
		return 0, err
	}

	return data[0], nil
}

func verifyId(d *i2c.I2C) bool{
	data, _, err := d.ReadRegBytes(CCS811_HW_ID, 1)
	if err == nil && data[0] == CCS811_HW_ID_CODE {
		return true
	}

	return false
}

func setConfig(d *i2c.I2C) error {
	data, _, err := d.ReadRegBytes(CCS811_STATUS, 1); if err != nil {
		return err
	}

	bin1:=0x01 & InterruptThreshold
	bin2:=0x01 & InterruptMode
	bin3:=0x03 & SamplingRate
	data[0]= bin1 << 2 | bin2 << 3 | bin3 << 4

	_, err = d.WriteBytes(append([]byte{CCS811_MEAS_MODE}, data...))

	return err
}

func setReset(d *i2c.I2C) error {
	_, err := d.WriteBytes([]byte {CCS811_SW_RESET, 0x11, 0xE5, 0x72, 0x8A})

	return err
}
