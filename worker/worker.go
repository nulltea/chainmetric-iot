package worker

import (
	"sync"
	"time"

	"sensorsys/model"
)

type ReceiverFunc func(model.MetricReadings)

type SensorsReader struct {
	ctx          *Context
	waitGroup    sync.WaitGroup
	sensors      []Sensor
	receivers    []ReceiverFunc
}

func NewSensorsReader(ctx *Context) *SensorsReader {
	return &SensorsReader{
		ctx:     ctx,
		sensors: make([]Sensor, 0),
	}
}

func (s *SensorsReader) RegisterSensor(sensor Sensor) {
	s.sensors = append(s.sensors, sensor)
}

func (s *SensorsReader) RegisterSensors(sensors ...Sensor) {
	for _, sensor := range sensors {
		s.RegisterSensor(sensor)
	}
}

func (s *SensorsReader) SubscribeReceiver(receiver ReceiverFunc, metrics ...model.Metric) {
	for _, metric := range metrics {
		if _, ok := s.ctx.Pipe[metric]; ok {
			continue
		}
		s.ctx.Pipe[metric] = make(chan model.MetricReading)
	}
	s.receivers = append(s.receivers, receiver)
}

func (s *SensorsReader) Process() {
	s.initSensors()

	go func() {
		for {
			for _, sensor := range s.sensors {
				go sensor.Harvest(s.ctx.ForSensor(sensor))
			}
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			for _, receiver := range s.receivers {
				go receiver(nil)
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func (s *SensorsReader) ProcessSync() model.MetricReadings {
	s.initSensors()
	// s.readingsPipe = make(chan model.MetricReadings, len(s.routines))
	// s.waitGroup.Add(len(s.routines))
	// go func() {
	// 	s.waitGroup.Wait()
	// 	close(s.readingsPipe)
	// }()
	// for _, routine := range s.routines {
	// 	go routine(s.ctx)
	// }
	readings := model.MetricReadings{}
	// for reading := range s.readingsPipe {
	// 	for metric, value := range reading {
	// 		readings[metric] = value
	// 	}
	// }
	return readings
}

func (s *SensorsReader) Clean() {
	for _, sensor := range s.sensors {
		if err := sensor.Close(); err != nil {
			s.ctx.ForSensor(sensor).Error(err)
		}
	}
}

func (s *SensorsReader) initSensors() {
	for i, sensor := range s.sensors {
		if err := sensor.Init(); err != nil {
			s.ctx.ForSensor(sensor).Error(err)
			s.unregisterSensor(i)
		}
	}
}

func (s *SensorsReader) unregisterSensor(i int) {
	copy(s.sensors[i:], s.sensors[i+1:])
	s.sensors = s.sensors[:len(s.sensors)-1]
}
