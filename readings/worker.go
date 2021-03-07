package readings

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/model"
	"github.com/timoth-y/iot-blockchain-sensorsys/readings/receiver"
	"github.com/timoth-y/iot-blockchain-sensorsys/readings/sensor"
)

type SensorsReader struct {
	context       *Context
	sensors       []sensor.Sensor
	requests      chan receiver.Request
	standbyTimers map[sensor.Sensor]*time.Timer
}

func NewSensorsReader(ctx *Context) *SensorsReader {
	return &SensorsReader{
		context:       ctx,
		sensors:       make([]sensor.Sensor, 0),
		requests:      make(chan receiver.Request),
		standbyTimers: make(map[sensor.Sensor]*time.Timer),
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

func (s *SensorsReader) SubscribeReceiver(handler receiver.ReceiverFunc, period time.Duration, metrics ...models.Metric) {
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

func (s *SensorsReader) SendRequest(handler receiver.ReceiverFunc, metrics ...models.Metric) {
	s.requests <- receiver.Request{
		Metrics: metrics,
		Handler: handler,
	}
}

func (s *SensorsReader) Process() {
	for {
		select {
		case req := <- s.requests:
			go s.handleRequest(s.context.ForRequest(req.Metrics), req)
		}
	}
}

func (s *SensorsReader) handleRequest(ctx *receiver.Context, req receiver.Request) {
	ctx.WaitGroup = &sync.WaitGroup{}

	for _, sn := range s.sensors {
		for _, metric := range req.Metrics {
			if suitable(sn, metric) {
				c := ctx.ForSensor(sn)

				if err := s.initSensor(sn); err != nil {
					c.Error(err)
					continue
				}

				ctx.WaitGroup.Add(1)

				c, cancel := c.SetTimeout(1 * time.Second) // TODO: configure or base on request period
				defer cancel()

				go s.readSensor(ctx, c, sn)

				break
			}
		}
	}

	ctx.WaitGroup.Wait()

	results := aggregate(ctx)
	req.Handler(results)

	return
}

func suitable(sensor sensor.Sensor, metric models.Metric) bool {
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

		if len(readings) != 0 { // TODO: config-based or precision-based aggregation here
			results[metric] = readings[len(readings) - 1].Value
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

func (s *SensorsReader) initSensor(sn sensor.Sensor) error {
	if !sn.Active() {
		if err := sn.Init(); err != nil {
			return err
		}
	}
	duration := time.Duration(s.context.Config.Worker.CloseOnStandbyTime) * time.Millisecond
	if timer, ok := s.standbyTimers[sn]; ok && timer != nil {
		if !timer.Reset(duration) {
			go handleStandby(timer, sn)
		}
	} else {
		s.standbyTimers[sn] = time.NewTimer(duration)
		go handleStandby(s.standbyTimers[sn], sn)
	}

	return nil
}

func (s *SensorsReader) readSensor(req *receiver.Context, ctx *sensor.Context, sn sensor.Sensor) {
	defer req.WaitGroup.Done()

	done := make(chan bool)

	go func() {
		sn.Harvest(ctx)
		done <- true
	}()

	select {
	case <- ctx.Done():
		switch ctx.Err() {
		case context.DeadlineExceeded:
			ctx.Error(errors.New("sensor reading timeout: time exceeded"))
		case context.Canceled:
			ctx.Info("sensor reading canceled by force")
		}
		return
	case <- done:
		return
	}
}

func handleStandby(t *time.Timer, sn sensor.Sensor) {
	<- t.C
	sn.Close()
}
