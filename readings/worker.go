package readings

import (
	"context"

	"sensor/model"
)

type sensorReadingRoutine func(ctx context.Context)

type SensorsReader struct {
	ctx context.Context
	routines []sensorReadingRoutine
	readings chan model.MetricReadings
}

func NewSensorsReader(ctx context.Context) *SensorsReader {
	return &SensorsReader{
		ctx: ctx,
		routines: []sensorReadingRoutine{},
	}
}

func (s *SensorsReader) Execute() model.MetricReadings {
	s.readings = make(chan model.MetricReadings, len(s.routines))
	for _, routine := range s.routines {
		go routine(s.ctx)
	}
	readings := model.MetricReadings{}
	for reading := range s.readings {
		for metric, value := range reading {
			readings[metric] = value
		}
	}
	return readings
}

func (s *SensorsReader) subscribe(routine sensorReadingRoutine) {
	s.routines = append(s.routines, routine)
}