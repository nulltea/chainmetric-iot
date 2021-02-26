package readings

import (
	"context"
	"fmt"

	"sensorsys/model"
	"sensorsys/sensors"
)

func (s *SensorsReader) SubscribeToAccelerationReadings(addr uint8, bus int) error {
	sensor := sensors.NewADXL345(addr, bus)
	s.deferQueue = append(s.deferQueue, sensor.Close)

	if err := sensor.Init(); err != nil {
		return err
	}

	s.subscribe(func(ctx context.Context) {
		defer s.waitGroup.Done()

		vector, err := sensor.Read(); if err != nil {
			fmt.Println(err)
			return
		}

		s.readings <- model.MetricReadings{
			model.Acceleration:   vector,
		}
	})

	return nil
}
