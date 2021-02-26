package readings

import (
	"context"
	"fmt"

	"sensorsys/model"
	"sensorsys/sensors"
)

func (s *SensorsReader) SubscribeToTemperatureReadings(deviceID string, pin int) error {
	sensor := sensors.NewDHTxx(deviceID, pin)

	s.subscribe(func(ctx context.Context) {
		defer s.waitGroup.Done()

		temperature, humidity, err := sensor.Read(); if err != nil {
			fmt.Println(err)
			return
		}

		s.readings <- model.MetricReadings{
			model.Temperature: temperature,
			model.Humidity: humidity,
		}
	})

	return nil
}

