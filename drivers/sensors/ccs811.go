package sensors

import (
	"fmt"
	"time"

	"github.com/d2r2/go-i2c"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
)

var (
	InterruptMode byte = 0
	InterruptThreshold byte = 0
	SamplingRate byte = CCS811_DRIVE_MODE_10SEC
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

func (s *CCS811) ID() string {
	return "CCS811"
}

func (s *CCS811) Init() (err error) {
	s. i2c, err = i2c.NewI2C(s.addr, s.bus); if err != nil {
		return
	}

	if !s.Verify() {
		return fmt.Errorf("not that sensorType")
	}

	s.setReset()
	time.Sleep(CCS811_RESET_TIME * time.Millisecond)

	s.getStatus()

	_, err = s.i2c.WriteBytes([]byte{CCS811_BOOTLOADER_APP_START}); if err != nil {
		return err
	}

	time.Sleep(CCS811_APP_START_TIME * time.Millisecond)

	status, err := s.getStatus(); if err != nil {
		return err
	}

	if status &CCS811_ERROR_BIT != 0 {
		return fmt.Errorf("CCS811 device has error")
	}

	if status &CCS811_FW_MODE_BIT == 0 {
		return fmt.Errorf("CCS811 device is in FW mode")
	}

	s.setConfig(); if err != nil {
		return err
	}

	return
}

func (s *CCS811) Read() (eCO2 float64, eTVOC float64, err error) {
	retry := 10
	for retry > 0 {
		retry--
		ready, err := s.isDataReady(); if err != nil {
			return 0, 0, err
		}
		if ready {
			buffer, _, err := s.i2c.ReadRegBytes(CCS811_ALG_RESULT_DATA, 4)
			if err != nil {
				return 0, 0, err
			}
			eCO2 = float64((uint16(buffer[0]) << 8) | uint16(buffer[1]))
			eTVOC = float64((uint16(buffer[2]) << 8) | uint16(buffer[3]))
			break
		}
		time.Sleep(CCS811_RETRY_TIME * time.Millisecond)
	}
	err = nil
	return
}

func (s *CCS811) Harvest(ctx *Context) {
	eCO2, eTVOC, err := s.Read()

	if eCO2 != 0 {
		ctx.For(metrics.AirCO2Concentration).Write(eCO2)
	}

	if eTVOC != 0 {
		ctx.For(metrics.AirTVOCsConcentration).Write(eTVOC)
	}

	ctx.Error(err)
}

func (s *CCS811) Metrics() []models.Metric {
	return []models.Metric {
		metrics.AirCO2Concentration,
		metrics.AirTVOCsConcentration,
	}
}

func (s *CCS811) Verify() bool {
	buffer, _, err := s.i2c.ReadRegBytes(CCS811_HW_ID, 1)
	if err == nil && buffer[0] == CCS811_HW_ID_CODE {
		return true
	}

	return false
}

func (s *CCS811) Active() bool {
	return s.i2c != nil
}

func (s *CCS811) Close() error {
	defer s.clean()
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

func (s *CCS811) clean() {
	s.i2c = nil
}
