package readings

import (
	"fmt"
	"sync"
	"time"

	"sensorsys/model"
	"sensorsys/readings/receiver"
	"sensorsys/readings/sensor"
)

type SensorsReader struct {
	context   *Context
	sensors   []sensor.Sensor
	requests  chan Request
}

func NewSensorsReader(ctx *Context) *SensorsReader {
	return &SensorsReader{
		context:   ctx,
		sensors:   make([]sensor.Sensor, 0),
		requests:  make(chan Request),
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
	go func() {
		for {
			s.requests <- Request {
				Metrics: metrics,
				Handler: receiver,
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
					go s.handle(s.context.ForRequest(req.Metrics), req)
			}
		}
	}()
}

func (s *SensorsReader) handle(ctx *receiver.Context, req Request) {
	ctx.WaitGroup = &sync.WaitGroup{}

	for _, metric := range req.Metrics {
		for _, sensor := range s.sensors {
			if suitable(sensor, metric) {
				ctx.WaitGroup.Add(1)

				go func() {
					sensor.Harvest(ctx.ForSensor(sensor))
					ctx.WaitGroup.Done()
				}()
			}
		}
	}

	ctx.WaitGroup.Wait()

	results := aggregate(ctx)
	req.Handler(results)

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

func aggregate(ctx *receiver.Context) model.MetricReadings {
	results := make(model.MetricReadings)
	for metric, ch := range ctx.Pipe {
		readings := make([]model.MetricReading, 0)
		L:
		for {
			select {
			case reading := <- ch:
				readings = append(readings, reading)
			default:
				break L
			}
		}
		if len(readings) != 0 {
			// TODO config-based or precision-based aggregation here

			results[metric] = readings[len(readings) - 1]
		} else {
			fmt.Println("empty")
		}
	}

	return results
}


func (s *SensorsReader) Clean() {
	for _, sensor := range s.sensors {
		if err := sensor.Close(); err != nil {
			s.context.ForSensor(sensor).Error(err)
		}
	}
}

func (s *SensorsReader) initSensors() {
	for i, sensor := range s.sensors {
		if err := sensor.Init(); err != nil {
			s.context.ForSensor(sensor).Error(err)
			s.unregisterSensor(i)
		}
	}
}

func (s *SensorsReader) unregisterSensor(i int) {
	copy(s.sensors[i:], s.sensors[i+1:])
	s.sensors = s.sensors[:len(s.sensors)-1]
}
