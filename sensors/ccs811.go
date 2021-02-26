package sensors

import (
	"fmt"
	"time"

	"github.com/d2r2/go-i2c"
)

const(
	// Registers
	CCS811_STATUS          = 0x00
	CCS811_MEAS_MODE       = 0x01
	CCS811_ALG_RESULT_DATA = 0x02
	CCS811_RAW_DATA        = 0x03
	CCS811_ENV_DATA        = 0x05
	CCS811_NTC             = 0x06
	CCS811_THRESHOLDS      = 0x10
	CCS811_BASELINE        = 0x11
	CCS811_HW_ID           = 0x20
	CCS811_HW_VERSION      = 0x21
	CCS811_FW_BOOT_VERSION = 0x23
	CCS811_FW_APP_VERSION  = 0x24
	CCS811_ERROR_ID        = 0xE0
	CCS811_SW_RESET        = 0xFF

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

	// Constants
	CCS811_HW_ID_CODE   = 0x81
	CCS811_REF_RESISTOR = 100000

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
	CCS811_RETRY_TIME = 1000
)

var (
	InterruptMode byte = 0
	InterruptThreshold byte = 0
	SamplingRate byte =  CCS811_DRIVE_MODE_10SEC
)

type CCS811 struct {
	addr uint8
	bus int
	i2c *i2c.I2C
}

func NewCCS811(addr uint8, bus int) *CCS811 {
	return &CCS811{
		addr: addr,
		bus: bus,
	}
}

func (s *CCS811) Init() (err error) {
	s. i2c, err = i2c.NewI2C(s.addr, s.bus); if err != nil {
		return
	}

	if !s.Verify() {
		return fmt.Errorf("not that sensorType")
	}

	// setReset(i2c)
	// time.Sleep(CCS811_RESET_TIME * time.Millisecond)

	s.getStatus()

	_, err = s.i2c.WriteBytes([]byte{CCS811_BOOTLOADER_APP_START}); if err != nil {
		return err
	}

	time.Sleep(CCS811_APP_START_TIME * time.Millisecond)

	status, err := s.getStatus(); if err != nil {
		return err
	}

	if status & CCS811_ERROR_BIT != 0 {
		return fmt.Errorf("CCS811 device has error")
	}

	if status & CCS811_FW_MODE_BIT == 0 {
		return fmt.Errorf("CCS811 device is in FW mode")
	}

	// setConfig(i2c); if err != nil {
	// 	return err
	// }

	return
}

func (s *CCS811) Read() (eCO2 float32, eTVOC float32, err error) {
	for {
		ready, err := s.isDataReady(); if err != nil {
			return 0, 0, err
		}
		if ready {
			buffer, _, err := s.i2c.ReadRegBytes(CCS811_ALG_RESULT_DATA, 4)
			if err != nil {
				return 0, 0, err
			}
			eCO2 = float32((uint16(buffer[0]) << 8) | uint16(buffer[1]))
			eTVOC = float32((uint16(buffer[2]) << 8) | uint16(buffer[3]))
			break
		}
		time.Sleep(CCS811_RETRY_TIME * time.Millisecond)
	}
	err = nil
	return
}

func (s *CCS811) Verify() bool {
	buffer, _, err := s.i2c.ReadRegBytes(CCS811_HW_ID, 1)
	if err == nil && buffer[0] == CCS811_HW_ID_CODE {
		return true
	}

	return false
}

func (s *CCS811) Close() error {
	return s.i2c.Close()
}

func (s *CCS811) isDataReady() (bool, error) {
	sts, err := s.getStatus()
	if err != nil {
		return false, err
	}

	return (sts & CCS811_DATA_READY_BIT) != 0, nil
}

func (s *CCS811) getStatus() (byte, error) {
	data, _, err := s.i2c.ReadRegBytes(CCS811_STATUS, 1); if err != nil {
		return 0, err
	}

	// fmt.Printf("status: %#x\n", data[0])

	return data[0], nil
}



func (s *CCS811) setConfig() error {
	buffer := make([]byte, 1)
	bin1 := 0x01 & InterruptThreshold
	bin2 := 0x01 & InterruptMode
	bin3 := 0x03 & SamplingRate
	buffer[0] = bin1 << 2 | bin2 << 3 | bin3 << 4

	_, err := s.i2c.WriteBytes(append([]byte{CCS811_MEAS_MODE}, buffer...))

	return err
}

func (s *CCS811) setReset() error {
	_, err := s.i2c.WriteBytes([]byte {CCS811_SW_RESET, 0x11, 0xE5, 0x72, 0x8A})

	return err
}
