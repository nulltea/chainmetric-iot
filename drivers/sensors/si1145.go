package sensors

import (
	"fmt"

	"github.com/d2r2/go-i2c"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
)

type SI1145 struct {
	addr uint8
	bus int
	i2c *i2c.I2C
}

func NewSI1145(addr uint8, bus int) *SI1145 {
	return &SI1145{
		addr: addr,
		bus: bus,
	}
}

func (s *SI1145) ID() string {
	return "SI1145"
}

func (s *SI1145) Init() (err error) {
	s.i2c, err = i2c.NewI2C(s.addr, s.bus); if err != nil {
		return
	}

	if !s.Verify() {
		return fmt.Errorf("not SI1145 sensorType")
	}

	// Enable UV index measurement coefficients
	s.i2c.WriteRegU8(SI1145_REG_UCOEFF0, 0x29)
	s.i2c.WriteRegU8(SI1145_REG_UCOEFF1, 0x89)
	s.i2c.WriteRegU8(SI1145_REG_UCOEFF2, 0x02)
	s.i2c.WriteRegU8(SI1145_REG_UCOEFF3, 0x00)

	// Enable UV sensorType
	s.writeParam(SI1145_PARAM_CHLIST,
		SI1145_PARAM_CHLIST_ENUV|SI1145_PARAM_CHLIST_ENALSIR|SI1145_PARAM_CHLIST_ENALSVIS|SI1145_PARAM_CHLIST_ENPS1)

	// Enable interrupt on every sample
	s.i2c.WriteRegU8(SI1145_REG_INTCFG, SI1145_REG_INTCFG_INTOE)
	s.i2c.WriteRegU8(SI1145_REG_IRQEN, SI1145_REG_IRQEN_ALSEVERYSAMPLE)

	// Program LED current
	s.i2c.WriteRegU8(SI1145_REG_PSLED21, 0x03) // 20mA for LED 1 only
	s.writeParam(SI1145_PARAM_PS1ADCMUX, SI1145_PARAM_ADCMUX_LARGEIR)

	// Proximity sensorType //1 uses LED //1
	s.writeParam(SI1145_PARAM_PSLED12SEL, SI1145_PARAM_PSLED12SEL_PS1LED1)

	// Fastest clocks, clock div 1
	s.writeParam(SI1145_PARAM_PSADCGAIN, 0)

	// Take 511 clocks to measure
	s.writeParam(SI1145_PARAM_PSADCOUNTER, SI1145_PARAM_ADCCOUNTER_511CLK)

	// in proximity mode, high range
	s.writeParam(SI1145_PARAM_PSADCMISC, SI1145_PARAM_PSADCMISC_RANGE|SI1145_PARAM_PSADCMISC_PSMODE)
	s.writeParam(SI1145_PARAM_ALSIRADCMUX, SI1145_PARAM_ADCMUX_SMALLIR)

	// Fastest clocks, clock div 1
	s.writeParam(SI1145_PARAM_ALSIRADCGAIN, 0)

	// Take 511 clocks to measure
	s.writeParam(SI1145_PARAM_ALSIRADCOUNTER, SI1145_PARAM_ADCCOUNTER_511CLK)

	// in high range mode
	s.writeParam(SI1145_PARAM_ALSIRADCMISC, SI1145_PARAM_ALSIRADCMISC_RANGE)

	// fastest clocks, clock div 1
	s.writeParam(SI1145_PARAM_ALSVISADCGAIN, 0)

	// Take 511 clocks to measure
	s.writeParam(SI1145_PARAM_ALSVISADCOUNTER, SI1145_PARAM_ADCCOUNTER_511CLK)

	// in high range mode (not normal signal)
	s.writeParam(SI1145_PARAM_ALSVISADCMISC, SI1145_PARAM_ALSVISADCMISC_VISRANGE)

	// measurement rate for auto
	s.i2c.WriteRegU8(SI1145_REG_MEASRATE0, 0xFF) // 255 * 31.25uS = 8ms

	// auto run
	s.i2c.WriteRegU8(SI1145_REG_COMMAND, SI1145_PSALS_AUTO)

	return nil
}

// ReadUV returns the UV index * 100 (divide by 100 to get the index)
func (s *SI1145) ReadUV() (float64, error) {
	res, err := s.i2c.ReadRegU16LE(SI1145_REG_UVINDEX0)
	return float64(res), err
}

// ReadVisible returns visible + IR light levels
func (s *SI1145) ReadVisible() (float64, error) {
	res, err := s.i2c.ReadRegU16LE(SI1145_REG_ALSVISDATA0)
	return float64(res), err
}

// ReadIR returns IR light levels
func (s *SI1145) ReadIR() (float64, error) {
	res, err := s.i2c.ReadRegU16LE(SI1145_REG_ALSIRDATA0)
	return float64(res), err
}

// ReadProximity returns "Proximity" - assumes an IR LED is attached to LED
func (s *SI1145) ReadProximity() (float64, error) {
	res, err := s.i2c.ReadRegU16LE(SI1145_REG_PS1DATA0)
	return float64(res), err
}

func (s *SI1145) Harvest(ctx *Context) {
	ctx.For(metrics.UVLight).WriteWithError(s.ReadUV())
	ctx.For(metrics.VisibleLight).WriteWithError(s.ReadVisible())
	ctx.For(metrics.IRLight).WriteWithError(s.ReadIR())
	ctx.For(metrics.Proximity).WriteWithError(s.ReadProximity())
}

func (s *SI1145) Metrics() []models.Metric {
	return []models.Metric {
		metrics.UVLight,
		metrics.VisibleLight,
		metrics.IRLight,
		metrics.Proximity,
	}
}

func (s *SI1145) Verify() bool {
	return true // TODO verify by ID
}

func (s *SI1145) Active() bool {
	return s.i2c != nil
}

func (s *SI1145) Close() error {
	defer s.clean()
	return s.i2c.Close()
}

func (s *SI1145) writeParam(p, v uint8) (uint8, error) {
	s.i2c.WriteRegU8(SI1145_REG_PARAMWR, v)
	s.i2c.WriteRegU8(SI1145_REG_COMMAND, p |SI1145_PARAM_SET)

	return s.i2c.ReadRegU8(SI1145_REG_PARAMRD)
}

func (s *SI1145) clean() {
	s.i2c = nil
}
