package readings

import (
	"context"
	"fmt"
	"sync"

	"sensorsys/model"
)

type sensorReadingRoutine func(ctx context.Context)

type SensorsReader struct {
	ctx context.Context
	waitGroup sync.WaitGroup
	routines []sensorReadingRoutine
	readings chan model.MetricReadings
	cleanQueue []func() error
}

func NewSensorsReader(ctx context.Context) *SensorsReader {
	return &SensorsReader{
		ctx: ctx,
		routines: []sensorReadingRoutine{},
		cleanQueue: []func() error {},
	}
}

func (s *SensorsReader) Execute() model.MetricReadings {
	s.readings = make(chan model.MetricReadings, len(s.routines))
	s.waitGroup.Add(len(s.routines))
	go func() {
		s.waitGroup.Wait()
		close(s.readings)
	}()
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

func (s *SensorsReader) Clean() {
	for _, cf := range s.cleanQueue {
		if err := cf(); err != nil {
			fmt.Println(err)
		}
	}
}
