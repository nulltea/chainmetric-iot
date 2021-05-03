package sensors

import (
	"math"

	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripheries"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
)

type MAX44009 struct {
	*peripheries.I2C
}

func NewMAX44009(addr uint16, bus int) sensor.Sensor {
	return &MAX44009{
		I2C: peripheries.NewI2C(addr, bus),
	}
}

func (s *MAX44009) ID() string {
	return "MAX44009"
}

func (s *MAX44009) Init() error {
	if err := s.I2C.Init(); err != nil {
		return err
	}

	return nil
}

func (s *MAX44009) Read() (float64, error) {
	buffer, err := s.ReadRegBytes(MAX44009_LUX_READING_REGISTER, 2); if err != nil {
		return math.NaN(), err
	}

	exponent := int((buffer[0] & 0xF0) >> 4)
	mantissa := int(((buffer[0] & 0x0F) << 4) | (buffer[1] & 0x0F))
	lux := math.Pow(float64(2), float64(exponent)) * float64(mantissa) * 0.045

	return lux, nil
}

func (s *MAX44009) Harvest(ctx *sensor.Context) {
	ctx.For(metrics.Luminosity).WriteWithError(s.Read())
}

func (s *MAX44009) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Luminosity,
	}
}

func (s *MAX44009) Verify() bool {
	if !s.I2C.Verify() {
		return false
	}

	if devID, err := s.I2C.ReadReg(MAX44009_DEVICE_ID_REGISTER); err == nil {
		return devID == MAX44009_DEVICE_ID
	}

	return true
}
