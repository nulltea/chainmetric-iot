package sensors

import (
	"sync"

	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-sensorsys/core"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-sensorsys/core/sensor"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/periphery"
)

var (
	si1145Mutex = &sync.Mutex{}
)

type SI1145 struct {
	*periphery.I2C
}

func NewSI1145(addr uint16, bus int) core.Sensor {
	return &SI1145{
		I2C: periphery.NewI2C(addr, bus, periphery.WithMutex(si1145Mutex)),
	}
}

func (s *SI1145) ID() string {
	return "SI1145"
}

func (s *SI1145) Init() (err error) {
	if err = s.I2C.Init(); err != nil {
		return
	}

	// Enable UV index measurement coefficients
	err = s.WriteRegBytes(SI1145_REG_UCOEFF0, 0x7B)
	err = s.WriteRegBytes(SI1145_REG_UCOEFF1, 0x6B)
	err = s.WriteRegBytes(SI1145_REG_UCOEFF2, 0x01)
	err = s.WriteRegBytes(SI1145_REG_UCOEFF3, 0x00)

	// Enable UV sensorType
	_, err = s.writeParam(SI1145_PARAM_CHLIST,
		SI1145_PARAM_CHLIST_ENUV    | SI1145_PARAM_CHLIST_ENAUX |
		SI1145_PARAM_CHLIST_ENALSIR | SI1145_PARAM_CHLIST_ENALSVIS | SI1145_PARAM_CHLIST_ENPS1)

	// Enable interrupt on every sample
	err = s.WriteRegBytes(SI1145_REG_INTCFG, SI1145_REG_INTCFG_INTOE)
	err = s.WriteRegBytes(SI1145_REG_IRQEN, SI1145_REG_IRQEN_ALSEVERYSAMPLE)

	// Program LED current
	err = s.WriteRegBytes(SI1145_REG_PSLED21, 0x03) // 20mA for LED 1 only
	_, err = s.writeParam(SI1145_PARAM_PS1ADCMUX, SI1145_PARAM_ADCMUX_LARGEIR)

	// Proximity sensorType //1 uses LED //1
	_, err = s.writeParam(SI1145_PARAM_PSLED12SEL, SI1145_PARAM_PSLED12SEL_PS1LED1)

	// Fastest clocks, clock div 1
	_, err = s.writeParam(SI1145_PARAM_PSADCGAIN, 0)

	// Take 511 clocks to measure
	_, err = s.writeParam(SI1145_PARAM_PSADCOUNTER, SI1145_PARAM_ADCCOUNTER_511CLK)

	// in proximity mode, high range
	_, err = s.writeParam(SI1145_PARAM_PSADCMISC, SI1145_PARAM_PSADCMISC_RANGE|SI1145_PARAM_PSADCMISC_PSMODE)
	_, err = s.writeParam(SI1145_PARAM_ALSIRADCMUX, SI1145_PARAM_ADCMUX_SMALLIR)

	// Fastest clocks, clock div 1
	_, err = s.writeParam(SI1145_PARAM_ALSIRADCGAIN, 0)

	// Take 511 clocks to measure
	_, err = s.writeParam(SI1145_PARAM_ALSIRADCOUNTER, SI1145_PARAM_ADCCOUNTER_511CLK)

	// in high range mode
	_, err = s.writeParam(SI1145_PARAM_ALSIRADCMISC, SI1145_PARAM_ALSIRADCMISC_RANGE)

	// fastest clocks, clock div 1
	_, err = s.writeParam(SI1145_PARAM_ALSVISADCGAIN, 0)

	// Take 511 clocks to measure
	_, err = s.writeParam(SI1145_PARAM_ALSVISADCOUNTER, SI1145_PARAM_ADCCOUNTER_511CLK)

	// in high range mode (not normal signal)
	_, err = s.writeParam(SI1145_PARAM_ALSVISADCMISC, SI1145_PARAM_ALSVISADCMISC_VISRANGE)

	// measurement rate for auto
	err = s.WriteRegBytes(SI1145_REG_MEASRATE0, 0xFF) // 255 * 31.25uS = 8ms

	// auto run
	err = s.WriteRegBytes(SI1145_REG_COMMAND, SI1145_PSALS_AUTO)

	return nil
}

// ReadUV returns the UV index * 100 (divide by 100 to get the index)
func (s *SI1145) ReadUV() (float64, error) {
	res, err := s.ReadRegU16LE(SI1145_REG_UVINDEX0)
	return float64(res), err
}

// ReadVisible returns visible + IR light levels
func (s *SI1145) ReadVisible() (float64, error) {
	res, err := s.ReadRegU16LE(SI1145_REG_ALSVISDATA0)
	return float64(res), err
}

// ReadIR returns IR light levels
func (s *SI1145) ReadIR() (float64, error) {
	res, err := s.ReadRegU16LE(SI1145_REG_ALSIRDATA0)
	return float64(res), err
}

// ReadProximity returns "Proximity" - assumes an IR LED is attached to LED
func (s *SI1145) ReadProximity() (float64, error) {
	res, err := s.ReadRegU16LE(SI1145_REG_PS1DATA0)
	return float64(res), err
}

func (s *SI1145) Harvest(ctx *sensor.Context) {
	max44009Mutex.Lock()
	defer max44009Mutex.Unlock()

	ctx.WriterFor(metrics.UVLight).WriteWithError(s.ReadUV())
	ctx.WriterFor(metrics.VisibleLight).WriteWithError(s.ReadVisible())
	ctx.WriterFor(metrics.IRLight).WriteWithError(s.ReadIR())
	ctx.WriterFor(metrics.Proximity).WriteWithError(s.ReadProximity())
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
	if !s.I2C.Verify() {
		return false
	}

	if devID, err := s.I2C.ReadReg(SI1145_DEVICE_ID_REGISTER); err == nil {
		return devID == SI1145_DEVICE_ID
	}

	return false
}

func (s *SI1145) writeParam(p, v uint8) (uint8, error) {
	if err := s.WriteRegBytes(SI1145_REG_PARAMWR, v); err != nil {
		return 0, err
	}

	if err := s.WriteRegBytes(SI1145_REG_COMMAND, p | SI1145_PARAM_SET); err != nil {
		return 0, err
	}

	return s.ReadReg(SI1145_REG_PARAMRD)
}
