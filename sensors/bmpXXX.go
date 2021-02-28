package sensors

import (
	"github.com/d2r2/go-bsbmp"
	"github.com/d2r2/go-i2c"

	"sensorsys/model/metrics"
	"sensorsys/worker"
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

func (s *BMPxx) Harvest(ctx *worker.Context) {
	ctx.For(metrics.Pressure).WriteWithError(s.bmp.ReadPressurePa(bsbmp.ACCURACY_STANDARD))
}

func (s *BMPxx) Close() error {
	return s.Close()
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
