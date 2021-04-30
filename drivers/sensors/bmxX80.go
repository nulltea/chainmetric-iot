package sensors

import (
	"math"

	"github.com/timoth-y/chainmetric-core/models"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/devices/bmxx80"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripheries"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
)

type BMP280 struct {
	*peripheries.I2C
	bmp  *bmxx80.Dev
}

func NewBMXX80(addr uint16, bus int) sensor.Sensor {
	return &BMP280{
		I2C: peripheries.NewI2C(addr, bus),
	}
}

func (s *BMP280) ID() string {
	return "BMP280"
}

func (s *BMP280) Init() (err error) {
	if err = s.I2C.Init(); err != nil {
		return
	}

	if s.bmp, err = bmxx80.NewI2C(s.Bus, s.Addr, &bmxx80.DefaultOpts); err != nil {
		return
	}

	return
}

func (s *BMP280) Harvest(ctx *sensor.Context) {
	var env = physic.Env{}

	if err := s.bmp.Sense(&env); err != nil {
		ctx.Error(err)

		return
	}

	ctx.For(metrics.Pressure).Write(float64(env.Pressure))
	ctx.For(metrics.Altitude).Write(s.pressureToAltitude(float64(env.Pressure)))
	ctx.For(metrics.Temperature).Write(env.Temperature.Celsius())
	// ctx.For(metrics.Humidity).Write(float64(env.Humidity)) TODO: test compatibility
}

func (s *BMP280) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Pressure,
		metrics.Altitude,
		metrics.Temperature,
		metrics.Humidity,
	}
}

func (s *BMP280) Verify() bool {
	if !s.I2C.Verify() {
		return false
	}

	if devID, err := s.I2C.ReadReg(BMP280_DEVICE_ID_REGISTER); err == nil {
		return devID == BMP280_DEVICE_ID
	}

	return false
}

func (s *BMP280) pressureToAltitude(p float64) float64 {
	// Approximate atmospheric pressure at sea level in Pa
	p0 := 1013250.0
	a := 44330 * (1 - math.Pow(p / p0, 1/5.255))
	// Round up to 2 decimals after point
	a2 := float64(int(a*100)) / 100
	return a2
}
