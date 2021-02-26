package readings

import (
	"context"
	"fmt"

	"sensorsys/model"
	"sensorsys/sensors"
)

func (s *SensorsReader) SubscribeToAirQualityReadings(addr uint8, bus int) error {
	sensor := sensors.NewCCS811(addr, bus)
	s.deferQueue = append(s.deferQueue, sensor.Close)

	if err := sensor.Init(); err != nil {
		return err
	}

	s.subscribe(func(ctx context.Context) {
		defer s.waitGroup.Done()

		eC02, eTVOC, err := sensor.Read(); if err != nil {
			fmt.Println(err)
			return
		}

		s.readings <- model.MetricReadings{
			model.AirCO2Concentration:   eC02,
			model.AirTVOCsConcentration: eTVOC,
		}
	})

	return nil
}
