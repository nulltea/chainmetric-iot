package sensors

import (
	"github.com/d2r2/go-bsbmp"
	"github.com/d2r2/go-i2c"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
)

type BMPxx struct {
	sensorType bsbmp.SensorType
	addr uint8
	bus int
	i2c *i2c.I2C
	bmp *bsbmp.BMP
}

func NewBMPxxx(deviceID string, addr uint8, bus int) *BMPxx {
	return &BMPxx{
		sensorType: sensorTypeBMP(deviceID),
	}
}

func NewBMP280(addr uint8, bus int) *BMPxx {
	return &BMPxx{
		sensorType: bsbmp.BME280,
	}
}

func (s *BMPxx) ID() string {
	return s.sensorType.String()
}

func (s *BMPxx) Init() (err error) {
	s.i2c, err = i2c.NewI2C(s.addr, s.bus); if err != nil {
		return
	}

	s.bmp, err = bsbmp.NewBMP(s.sensorType, s.i2c); if err != nil {
		return
	}

	return
}

func (s *BMPxx) Harvest(ctx *Context) {
	ctx.For(metrics.Pressure).WriteWithError(s.bmp.ReadPressurePa(bsbmp.ACCURACY_STANDARD))
}

func (s *BMPxx) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Pressure,
	}
}

func (s *BMPxx) Active() bool {
	return s.i2c != nil && s.bmp != nil
}

func (s *BMPxx) Close() error {
	defer s.clean()
	return s.i2c.Close()
}

func sensorTypeBMP(deviceID string) bsbmp.SensorType {
	switch deviceID {
	case "BMP180":
		return bsbmp.BMP280
	case "BMP280":
		return bsbmp.BMP280
	case "BME280":
		return bsbmp.BME280
	default:
		return bsbmp.BMP280
	}
}

func (s *BMPxx) clean() {
	s.i2c = nil
	s.bmp = nil
}
