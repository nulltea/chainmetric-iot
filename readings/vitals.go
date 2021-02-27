package readings

import (
	"context"
	"fmt"

	"sensorsys/model"
	"sensorsys/sensors"
)

func (s *SensorsReader) SubscribeToVitalsReadings( addr uint8, bus int) error {
	sensor := sensors.NewMAX30102(addr, bus)

	if err := sensor.Init(); err != nil {
		return err
	}
	s.deferQueue = append(s.deferQueue, sensor.Close)

	s.subscribe(func(ctx context.Context) {
		defer s.waitGroup.Done()

		heart, oxidation, err := sensor.Read(); if err != nil {
			fmt.Println(err)
			return
		}

		s.readings <- model.MetricReadings{
			model.HeartRate: heart,
			model.BloodOxidation: oxidation,
		}
	})

	return nil
}
