package sensors

import (
	"math"

	"github.com/d2r2/go-bsbmp"
	"github.com/timoth-y/iot-blockchain-contracts/models"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/devices/bmxx80"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/peripherals"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
)

type BMPxx struct {
	addr uint8
	bus  int
	dev  *peripherals.I2C
	bmp  *bmxx80.Dev
}

func NewBMXX80(addr uint16, bus int) *BMPxx {
	return &BMPxx{
		dev: peripherals.NewI2C(addr, bus),
	}
}

func (s *BMPxx) ID() string {
	return s.bmp.String()
}

func (s *BMPxx) Init() (err error) {
	if err = s.dev.Init(); err != nil {
		return
	}

	if s.bmp, err = bmxx80.NewI2C(s.dev.Bus, s.dev.Addr, &bmxx80.DefaultOpts); err != nil {
		return
	}

	s.bmp.String()

	return
}

func (s *BMPxx) Harvest(ctx *Context) {
	var env = physic.Env{}

	if err := s.bmp.Sense(&env); err != nil {
		ctx.Error(err)

		return
	}

	ctx.For(metrics.Pressure).Write(float64(env.Pressure))
	ctx.For(metrics.Altitude).Write(s.pressureToAltitude(float64(env.Pressure)))
	ctx.For(metrics.Temperature).Write(env.Temperature.Celsius())
	ctx.For(metrics.Humidity).Write(float64(env.Humidity))
}

func (s *BMPxx) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Pressure,
		metrics.Altitude,
		metrics.Temperature,
		metrics.Humidity,
	}
}

func (s *BMPxx) Active() bool {
	return s.dev != nil && s.bmp != nil
}

func (s *BMPxx) Close() error {
	defer s.clean()
	return s.dev.Close()
}

func sensorTypeBMP(deviceID string) bsbmp.SensorType {
	switch deviceID {
	case "BMP180":
		return bsbmp.BMP180
	case "BMP280":
		return bsbmp.BMP280
	case "BME280":
		return bsbmp.BME280
	default:
		return bsbmp.BMP280
	}
}

func (s *BMPxx) clean() {
	s.dev = nil
	s.bmp = nil
}

func (s *BMPxx) pressureToAltitude(p float64) float64 {
	// Approximate atmospheric pressure at sea level in Pa
	p0 := 1013250.0
	a := 44330 * (1 - math.Pow(p / p0, 1/5.255))
	// Round up to 2 decimals after point
	a2 := float64(int(a*100)) / 100
	return a2
}
