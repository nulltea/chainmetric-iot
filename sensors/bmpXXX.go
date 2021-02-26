package sensors

import (
	"github.com/d2r2/go-bsbmp"
	"github.com/d2r2/go-i2c"
)

type BMPxxx struct {
	sensorType bsbmp.SensorType
	addr uint8
	bus int
	i2c *i2c.I2C
	bmp *bsbmp.BMP
}

func NewBMPxxx(deviceID string, addr uint8, bus int) *BMPxxx {
	return &BMPxxx{
		sensorType: sensorTypeBMP(deviceID),
	}
}

func (s *BMPxxx) Init() (err error) {
	s.i2c, err = i2c.NewI2C(s.addr, s.bus); if err != nil {
		return
	}

	s.bmp, err = bsbmp.NewBMP(s.sensorType, s.i2c); if err != nil {
		return
	}

	return
}

func (s *BMPxxx) Read() (pa float32, err error) {
	return s.bmp.ReadPressurePa(bsbmp.ACCURACY_STANDARD)
}

func (s *BMPxxx) Close() error {
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
