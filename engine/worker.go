package engine

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/sensors"
	"github.com/timoth-y/iot-blockchain-sensorsys/model"
)

type SensorsReader struct {
	context       *Context
	sensors       []sensors.Sensor
	requests      chan Request
	standbyTimers map[sensors.Sensor]*time.Timer
	done          chan struct{}
}

func NewSensorsReader() *SensorsReader {
	return &SensorsReader{
		sensors:       make([]sensors.Sensor, 0),
		requests:      make(chan Request),
		standbyTimers: make(map[sensors.Sensor]*time.Timer),
		done:          make(chan struct{}),
	}
}

func (s *SensorsReader) Init(ctx *Context) error {
	s.context = ctx
	return nil
}

func (s *SensorsReader) RegisterSensors(sensors ...sensors.Sensor) {
	s.sensors = append(s.sensors, sensors...)
}

func (s *SensorsReader) SubscribeReceiver(handler ReceiverFunc, period time.Duration, metrics ...models.Metric) context.CancelFunc {
	ctx, cancel := context.WithCancel(s.context)
	go func() {
		for {
			s.requests <- Request{
				Metrics: metrics,
				Handler: handler,
			}

			select {
			case <- ctx.Done():
				return
			default:
				time.Sleep(period)
			}
		}
	}()

	return cancel
}

func (s *SensorsReader) SendRequest(handler ReceiverFunc, metrics ...models.Metric) {
	s.requests <- Request{
		Metrics: metrics,
		Handler: handler,
	}
}

func (s *SensorsReader) Process() {
	for {
		select {
		case request := <- s.requests:
			go s.handleRequest(request)
		case <- s.context.Done():
			return
		case <- s.done:
			s.context.Info("Reader process ended")
			return
		}
	}
}

func (s *SensorsReader) handleRequest(req Request) {
	var (
		waitGroup = &sync.WaitGroup{}
		pipe = make(model.MetricReadingsPipe)
	)


	for _, metric := range req.Metrics {
		pipe[metric] = make(chan model.MetricReading, 3)
	}

	for _, sn := range s.sensors {
		for _, metric := range req.Metrics {
			if suitable(sn, metric) {
				ctx := s.context.ForSensor(sn)
				ctx.Pipe = pipe

				if err := s.initSensor(sn); err != nil {
					ctx.Error(err)
					continue
				}

				waitGroup.Add(1)

				ctx, cancel := ctx.SetTimeout(2 * time.Second) // TODO: configure or base on request period
				defer cancel()

				go s.readSensor(ctx, sn, waitGroup)

				break
			}
		}
	}

	waitGroup.Wait()

	results := aggregate(pipe)
	req.Handler(results)

	return
}

func suitable(sensor sensors.Sensor, metric models.Metric) bool {
	for _, m := range sensor.Metrics() {
		if metric == m {
			return true
		}
	}

	return false
}

func aggregate(pipe model.MetricReadingsPipe) model.MetricReadings {
	results := make(model.MetricReadings)
	for metric, ch := range pipe {
		readings := make([]model.MetricReading, 0)

	LOOP: for {
			select {
			case reading := <- ch:
				readings = append(readings, reading)
			default:
				break LOOP
			}
		}

		if len(readings) != 0 {
			results[metric] = selectResult(readings)
		}
	}

	return results
}


func (s *SensorsReader) Close() {
	s.done <- struct{}{}
	for _, sensor := range s.sensors {
		if sensor.Active() {
			if err := sensor.Close(); err != nil {
				s.context.ForSensor(sensor).Error(err)
			}
		}
	}
}

func (s *SensorsReader) initSensor(sn sensors.Sensor) error {
	var (
		standby = viper.GetDuration("engine.sensor_sleep_standby_timeout")
	)

	if !sn.Active() {
		if err := sn.Init(); err != nil {
			return err
		}
	}

	if timer, ok := s.standbyTimers[sn]; ok && timer != nil {
		if !timer.Reset(standby) {
			go handleStandby(timer, sn)
		}
	} else {
		s.standbyTimers[sn] = time.NewTimer(standby)
		go handleStandby(s.standbyTimers[sn], sn)
	}

	return nil
}

func (s *SensorsReader) readSensor(ctx *sensors.Context, sn sensors.Sensor, wg *sync.WaitGroup) {
	defer wg.Done()

	if !sn.Active() {
		ctx.Warning("attempt of reading from non-active sensor")

		return
	}

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

func handleStandby(t *time.Timer, sn sensors.Sensor) {
	<- t.C
	sn.Close()
}

func selectResult(results []model.MetricReading) (result float64) {
	var (
		getPrecision = func(v float64) int {
			s := strconv.FormatFloat(v, 'f', -1, 64)
			i := strings.IndexByte(s, '.')
			if i > -1 {
				return len(s) - i - 1
			}
			return 0
		}
	)
	result = results[0].Value
	lastPrecision := 0

	for i := range results {
		precision := getPrecision(results[i].Value)

		if precision > lastPrecision {
			result = results[i].Value
		}

		lastPrecision = precision
	}

	return
}
