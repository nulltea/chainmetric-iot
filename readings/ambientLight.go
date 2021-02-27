package readings

import (
	"context"
	"fmt"

	"sensorsys/model"
	"sensorsys/sensors"
)

func (s *SensorsReader) SubscribeToAmbientLightReadings(addr uint8, bus int) error {
	sensor := sensors.NewSI1145(addr, bus)

	err := sensor.Init(); if err != nil {
		return err
	}
	s.deferQueue = append(s.deferQueue, sensor.Close)

	s.subscribe(func(ctx context.Context) {
		defer s.waitGroup.Done()

		uv, visible, ir, err := sensor.Read(); if err != nil {
			fmt.Println(err)
			return
		}

		s.readings <- model.MetricReadings{
			model.UVLight: uv,
			model.VisibleLight: visible,
			model.IRLight: ir,
		}
	})

	return nil
}
