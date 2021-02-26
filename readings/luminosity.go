package readings

import (
	"context"
	"fmt"

	"sensorsys/model"
	"sensorsys/sensors"
)

func (s *SensorsReader) SubscribeToLuminosityReadings(addr uint8, bus int) error {
	sensor := sensors.NewMAX44009(addr, bus)
	s.deferQueue = append(s.deferQueue, sensor.Close)

	if err := sensor.Init(); err != nil {
		return err
	}

	s.subscribe(func(ctx context.Context) {
		defer s.waitGroup.Done()

		lumen, err := sensor.Read(); if err != nil {
			fmt.Println(err)
			return
		}

		s.readings <- model.MetricReadings{
			model.Luminosity: lumen,
		}
	})

	return nil
}
