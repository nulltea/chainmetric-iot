package readings

import (
	"sync"
	"time"

	"sensorsys/model"
	"sensorsys/readings/receiver"
	"sensorsys/readings/sensor"
)

type SensorsReader struct {
	ctx       *Context
	waitGroup sync.WaitGroup
	sensors   []sensor.Sensor
	receivers []Receiver
	requests  chan Request
}

func NewSensorsReader(ctx *Context) *SensorsReader {
	return &SensorsReader{
		ctx:     ctx,
		sensors: make([]sensor.Sensor, 0),
		receivers: make([]Receiver, 0),
		requests: make(chan Request),
	}
}

func (s *SensorsReader) RegisterSensor(sensor sensor.Sensor) {
	s.sensors = append(s.sensors, sensor)
}

func (s *SensorsReader) RegisterSensors(sensors ...sensor.Sensor) {
	for _, sensor := range sensors {
		s.RegisterSensor(sensor)
	}
}

func (s *SensorsReader) SubscribeReceiver(receiver ReceiverFunc, period time.Duration, metrics ...model.Metric) {
	s.receivers = append(s.receivers, Receiver {
		Handler: receiver,
		Metrics: metrics,
		Period: period,
	})

	go func() {
		for {
			s.requests <- Request {
				Context: s.ctx.ForReceiver(),
				Metrics: metrics,
			}
			time.Sleep(period)
		}
	}()
}

func (s *SensorsReader) Process() {
	s.initSensors()
	go func() {
		for {
			// wait before receiver request
			select {
				case req := <- s.requests:
					go s.handle(req.Context, req.Metrics)
			}
		}
	}()
}

func (s *SensorsReader) handle(ctx *receiver.Context, metrics []model.Metric) {
	for _, metric := range metrics {
		for _, sensor := range s.sensors {
			if suitable(sensor, metric) {
				go sensor.Harvest(s.ctx.ForSensor(sensor))
			}
		}
	}
	return
}

func suitable(sensor sensor.Sensor, metric model.Metric) bool {
	for _, m := range sensor.Metrics() {
		if metric == m {
			return true
		}
	}
	return false
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

func (s *SensorsReader) Aggregate() {
	results := make(model.MetricReadings)
	for metric, ch := range s.ctx.Pipe {
		readings := make([]model.MetricReading, 0)
		for {
			select {
			case reading := <- ch:
				readings = append(readings, reading)
			default:
				break
			}
		}
		if len(readings) != 0 {
			// TODO config-based or precision-based aggregation here

			results[metric] = readings[len(readings) - 1]
		}
	}
	return
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
