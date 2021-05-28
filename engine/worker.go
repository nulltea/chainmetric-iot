package engine

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
	"github.com/timoth-y/chainmetric-sensorsys/model"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

type SensorsReader struct {
	context       *Context
	sensors       map[string]sensor.Sensor
	requests      chan Request
	standbyTimers map[sensor.Sensor]*time.Timer
	done          chan struct{}
}

func NewSensorsReader() *SensorsReader {
	return &SensorsReader{
		context:       NewContext(context.Background()),
		sensors:       make(map[string]sensor.Sensor),
		requests:      make(chan Request),
		standbyTimers: make(map[sensor.Sensor]*time.Timer),
		done:          make(chan struct{}, 1),
	}
}

func (r *SensorsReader) RegisterSensors(sensors ...sensor.Sensor) {
	for i, sensor := range sensors {
		r.sensors[sensor.ID()] = sensors[i]
	}
}

func (r *SensorsReader) UnregisterSensor(id string) {
	if sensor, ok := r.sensors[id]; ok {
		if sensor.Active() {
			sensor.Close()
		}
		delete(r.sensors, id)
	}
}

func (r *SensorsReader) RegisteredSensors() sensor.SensorsRegister {
	sMap := make(map[string]sensor.Sensor)

	for i, sensor := range r.sensors {
		sMap[sensor.ID()] = r.sensors[i]
	}

	return sMap
}

func (r *SensorsReader) SubscribeReceiver(handler ReceiverFunc, period time.Duration, metrics ...models.Metric) context.CancelFunc {
	ctx, cancel := context.WithCancel(r.context)
	go func() {
		for {
			r.requests <- Request{
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

func (r *SensorsReader) SendRequest(handler ReceiverFunc, metrics ...models.Metric) {
	r.requests <- Request{
		Metrics: metrics,
		Handler: handler,
	}
}

func (r *SensorsReader) Process(ctx context.Context) {
	for {
		select {
		case request := <- r.requests:
			go r.handleRequest(request)
		case <- ctx.Done():
			return
		case <- r.context.Done():
			return
		case <- r.done:
			shared.Logger.Debug("Sensors reader process ended.")
			return // TODO: refactor (too many contexts)
		}
	}
}

func (r *SensorsReader) handleRequest(req Request) {
	var (
		waitGroup = &sync.WaitGroup{}
		pipe = make(model.SensorReadingsPipe)
	)


	for _, metric := range req.Metrics {
		pipe[metric] = make(chan model.SensorReading, 3)
	}

	for _, sn := range r.sensors {
		for _, metric := range req.Metrics {
			if suitable(sn, metric) {
				ctx := r.context.ForSensor(sn)
				ctx.Pipe = pipe

				if err := r.initSensor(sn); err != nil {
					ctx.Error(err)
					continue
				}

				waitGroup.Add(1)

				ctx, cancel := ctx.SetTimeout(2 * time.Second) // TODO: configure or base on request period
				defer cancel()

				go r.readSensor(ctx, sn, waitGroup)

				break
			}
		}
	}

	waitGroup.Wait()

	results := aggregate(pipe)
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

func aggregate(pipe model.SensorReadingsPipe) model.SensorsReadingResults {
	results := make(model.SensorsReadingResults)
	for metric, ch := range pipe {
		readings := make([]model.SensorReading, 0)

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


func (r *SensorsReader) Close() {
	close(r.done)

	for _, sensor := range r.sensors {
		if sensor.Active() {
			if err := sensor.Close(); err != nil {
				r.context.ForSensor(sensor).Error(err)
			}
		}
	}
}


func (r *SensorsReader) initSensor(sn sensor.Sensor) error {
	var (
		standby = viper.GetDuration("engine.sensor_sleep_standby_timeout")
	)

	if !sn.Active() {
		if err := sn.Init(); err != nil {
			return err
		}
	}

	if timer, ok := r.standbyTimers[sn]; ok && timer != nil {
		if !timer.Reset(standby) {
			go handleStandby(timer, sn)
		}
	} else {
		r.standbyTimers[sn] = time.NewTimer(standby)
		go handleStandby(r.standbyTimers[sn], sn)
	}

	return nil
}

func (r *SensorsReader) readSensor(ctx *sensor.Context, sn sensor.Sensor, wg *sync.WaitGroup) {
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

func handleStandby(t *time.Timer, sn sensor.Sensor) {
	<-t.C
	sn.Close()
}

func selectResult(results []model.SensorReading) (result float64) {
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
