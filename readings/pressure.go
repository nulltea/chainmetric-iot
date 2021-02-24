package readings

import (
	"context"
	"fmt"

	"github.com/d2r2/go-bsbmp"
	"github.com/d2r2/go-i2c"

	"sensorsys/model"
)

func (s *SensorsReader) SubscribeToPressureReadings(sensor string, addr uint8, bus int) error {
	s.subscribe(func(ctx context.Context) {
		defer s.waitGroup.Done()
		pressure, err := readBMP(sensorTypeBMP(sensor), addr, bus); if err != nil {
			fmt.Println(err)
			return
		}
		s.readings <- model.MetricReadings{
			model.Pressure: pressure,
		}
	})

	return nil
}

func readBMP(stype bsbmp.SensorType, addr uint8, bus int) (float32, error) {
	i2c, err := i2c.NewI2C(addr, bus); if err != nil {
		return 0, err
	}
	defer i2c.Close()

	sensor, err := bsbmp.NewBMP(stype, i2c); if err != nil {
		return 0, err
	}

	return sensor.ReadPressurePa(bsbmp.ACCURACY_STANDARD)
}

func sensorTypeBMP(sensor string) bsbmp.SensorType {
	switch sensor {
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
