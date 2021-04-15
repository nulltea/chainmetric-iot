package sensors

import (
	"fmt"
	"math"

	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/peripherals"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
)

type MAX44009 struct {
	dev *peripherals.I2C
}

func NewMAX44009(addr uint16, bus int) *MAX44009 {
	return &MAX44009{
		dev: peripherals.NewI2C(addr, bus),
	}
}

func (s *MAX44009) ID() string {
	return "MAX44009"
}

func (s *MAX44009) Init() error {
	if err := s.dev.Init(); err != nil {
		return err
	}

	if !s.Verify() {
		return fmt.Errorf("driver is not compatiple with specified sensor")
	}

	if err := s.dev.WriteBytes(MAX44009_APP_START); err != nil {
		return err
	}

	return nil
}

func (s *MAX44009) Read() (float64, error) {
	buffer, err := s.dev.ReadBytes(2); if err != nil {
		return math.NaN(), err
	}

	return dataToLuminance(buffer), nil
}

func (s *MAX44009) Harvest(ctx *Context) {
	ctx.For(metrics.Luminosity).WriteWithError(s.Read())
}

func (s *MAX44009) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Luminosity,
	}
}

func (s *MAX44009) Verify() bool {
	return true
}

func (s *MAX44009) Active() bool {
	return s.dev.Active()
}

// Close disconnects from the device
func (s *MAX44009) Close() error {
	return s.dev.Close()
}

func dataToLuminance(d []byte) float64 {
	exponent := int((d[0] & 0xF0) >> 4)
	mantissa := int(((d[0] & 0x0F) << 4) | (d[1] & 0x0F))
	return math.Pow(float64(2), float64(exponent)) * float64(mantissa) * 0.045
}
