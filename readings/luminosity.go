package readings

import (
	"context"
	"fmt"
	"math"

	"github.com/d2r2/go-i2c"

	"sensorsys/model"
)

func (s *SensorsReader) SubscribeToLuminosityReadings(addr uint8, bus int) error {
	i2c, err := i2c.NewI2C(addr, bus); if err != nil {
		return err
	}
	s.cleanQueue = append(s.cleanQueue, i2c.Close)

	err = initMAX44009(i2c); if err != nil {
		return err
	}

	s.subscribe(func(ctx context.Context) {
		lumen, err := readMAX44009(i2c); if err != nil {
			fmt.Println(err)
		}
		s.readings <- model.MetricReadings{
			model.Luminosity: lumen,
		}
		s.waitGroup.Done()
	})

	return nil
}

func initMAX44009(i2c *i2c.I2C) error {
	_, err := i2c.WriteBytes([]byte{0x03}); if err != nil {
		return err
	}
	return nil
}

func readMAX44009(i2c *i2c.I2C) (float64, error) {
	_, err := i2c.WriteBytes([]byte{0x03}); if err != nil {
		return math.NaN(),err
	}

	var data = make([]byte, 2)
	_, err = i2c.ReadBytes(data); if err != nil {
		return math.NaN(),err
	}

	return dataToLuminance(data),nil
}

func dataToLuminance(d []byte) float64 {
	exponent := int((d[0] & 0xF0) >> 4)
	mantissa := int(((d[0] & 0x0F) << 4) | (d[1] & 0x0F))
	return math.Pow(float64(2), float64(exponent)) * float64(mantissa) * 0.045
}
