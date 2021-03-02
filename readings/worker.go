package readings

import (
	"sync"
	"time"

	"sensorsys/model"
	"sensorsys/readings/receiver"
	"sensorsys/readings/sensor"
)

type SensorsReader struct {
	context   *Context
	sensors   []sensor.Sensor
	requests  chan receiver.Request
}

func NewSensorsReader(ctx *Context) *SensorsReader {
	return &SensorsReader{
		context:   ctx,
		sensors:   make([]sensor.Sensor, 0),
		requests:  make(chan receiver.Request),
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

func (s *SensorsReader) SubscribeReceiver(handler receiver.ReceiverFunc, period time.Duration, metrics ...model.Metric) {
	go func() {
		for {
			s.requests <- receiver.Request{
				Metrics: metrics,
				Handler: handler,
			}

			time.Sleep(period)
		}
	}()
}

func (s *SensorsReader) Process() {
	go func() {
		for {
			select {
			case req := <- s.requests:
				go s.handle(s.context.ForRequest(req.Metrics), req)
			}
		}
	}()
}

func (s *SensorsReader) handle(ctx *receiver.Context, req receiver.Request) {
	ctx.WaitGroup = &sync.WaitGroup{}

	for _, sn := range s.sensors {
		for _, metric := range req.Metrics {
			if suitable(sn, metric) {
				snCtx := ctx.ForSensor(sn)

				if !sn.Active() {
					if err := sn.Init(); err != nil {
						snCtx.Error(err)
						continue
					}
				}

				ctx.WaitGroup.Add(1)

				go func(snCtx *sensor.Context, sensor sensor.Sensor) {
					sensor.Harvest(snCtx)
					ctx.WaitGroup.Done()
				}(snCtx, sn)
				break
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
			results[metric] = readings[len(readings) - 1].Value // TODO: config-based or precision-based aggregation here
		}
	}

	return results
}


func (s *SensorsReader) Clean() {
	for _, sensor := range s.sensors {
		if sensor.Active() {
			if err := sensor.Close(); err != nil {
				s.context.ForSensor(sensor).Error(err)
			}
		}
	}
}
