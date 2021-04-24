package sensors

import (
	"math"

	"github.com/timoth-y/iot-blockchain-contracts/models"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/devices/bmxx80"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/peripherals"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/sensor"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
)

type BMPxx struct {
	*peripherals.I2C
	bmp  *bmxx80.Dev
}

func NewBMXX80(addr uint16, bus int) sensor.Sensor {
	return &BMPxx{
		I2C: peripherals.NewI2C(addr, bus),
	}
}

func (s *BMPxx) ID() string {
	if s.bmp == nil {
		return "BMP280"
	}

	return s.bmp.String()
}

func (s *BMPxx) Init() (err error) {
	if err = s.I2C.Init(); err != nil {
		return
	}

	if s.bmp, err = bmxx80.NewI2C(s.Bus, s.Addr, &bmxx80.DefaultOpts); err != nil {
		return
	}

	return
}

func (s *BMPxx) Harvest(ctx *sensor.Context) {
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

func (s *BMPxx) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Pressure,
		metrics.Altitude,
		metrics.Temperature,
		metrics.Humidity,
	}
}

func (s *BMPxx) pressureToAltitude(p float64) float64 {
	// Approximate atmospheric pressure at sea level in Pa
	p0 := 1013250.0
	a := 44330 * (1 - math.Pow(p / p0, 1/5.255))
	// Round up to 2 decimals after point
	a2 := float64(int(a*100)) / 100
	return a2
}
