package sensors

import (
	"fmt"

	"github.com/d2r2/go-i2c"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/engine/sensor"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
)

const (
	// I2C
	HDC1080_ADDRESS =                       0x40
	// Registers
	HDC1080_TEMPERATURE_REGISTER =          0x00
	HDC1080_HUMIDITY_REGISTER =             0x01
	HDC1080_CONFIGURATION_REGISTER =        0x02

	// Configuration Register Bits
	HDC1080_CONFIG_RESET_BIT =             0x8000
	HDC1080_CONFIG_HEATER_ENABLE =          0x2000
	HDC1080_CONFIG_ACQUISITION_MODE =       0x1000
	HDC1080_CONFIG_BATTERY_STATUS =         0x0800
	HDC1080_CONFIG_TEMPERATURE_RESOLUTION = 0x0400
	HDC1080_CONFIG_HUMIDITY_RESOLUTION_HBIT =    0x0200
	HDC1080_CONFIG_HUMIDITY_RESOLUTION_LBIT =    0x0100
)

type HDC1080 struct {
	addr uint8
	bus int
	i2c *i2c.I2C
}

func NewHDC1080(addr uint8, bus int) *HDC1080 {
	return &HDC1080{
		addr: addr,
		bus: bus,
	}
}

func (s *HDC1080) ID() string {
	return "HDC1080"
}

func (s *HDC1080) Init() (err error) {
	s. i2c, err = i2c.NewI2C(s.addr, s.bus); if err != nil {
		return
	}

	if !s.Verify() {
		return fmt.Errorf("not HDC1080 sensor")
	}

	return
}

func (s *HDC1080) ReadTemperature() (float32, error) {
	data, _, err := s.i2c.ReadRegBytes(HDC1080_TEMPERATURE_REGISTER, 2); if err != nil {
		return 0, err
	}

	raw := float32(data[0] << 8 | data[1])

	return (raw / 65536.0) * 165.0 - 40.0, nil
}

func (s *HDC1080) ReadHumidity() (float32, error) {
	data, _, err := s.i2c.ReadRegBytes(HDC1080_HUMIDITY_REGISTER, 2); if err != nil {
		return 0, err
	}

	raw := float32(data[0] << 8 | data[1])

	return (raw / 65536.0) * 100.0, nil
}

func (s *HDC1080) Harvest(ctx *sensor.Context) {
	ctx.For(metrics.Temperature).WriteWithError(s.ReadTemperature())
	ctx.For(metrics.Humidity).WriteWithError(s.ReadHumidity())
}

func (s *HDC1080) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Temperature,
		metrics.Humidity,
	}
}

func (s *HDC1080) Verify() bool {
	return true
}

func (s *HDC1080) Active() bool {
	return s.i2c != nil
}

func (s *HDC1080) Close() error {
	defer s.clean()
	return s.i2c.Close()
}

func (s *HDC1080) clean() {
	s.i2c = nil
}
