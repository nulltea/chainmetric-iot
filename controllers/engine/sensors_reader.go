package engine

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-sensorsys/core"

	"github.com/timoth-y/chainmetric-sensorsys/core/sensor"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

type (
	// SensorsReader defines operational structure of the sensors reading engine.
	SensorsReader struct {
		once          *sync.Once
		sensors       sensor.SensorsRegister
		requests      chan request
		standbyTimers map[core.Sensor]*time.Timer
		active        bool
		cancel        context.CancelFunc
	}

	// ReadingResults defines map of values collected from core.Sensor for requested models.Metrics.
	ReadingResults map[models.Metric] float64

	// ReceiverFunc defines signature for sensor readings results receiver handler function.
	ReceiverFunc func(ReadingResults)

	// request stores data for the sensor readings request payload.
	request struct {
		Context context.Context
		Metrics []models.Metric
		Handler ReceiverFunc
	}
)

// NewSensorsReader constructs new SensorsReader instance.
func NewSensorsReader() *SensorsReader {
	return &SensorsReader{
		once:          &sync.Once{},
		sensors:       make(map[string]core.Sensor),
		requests:      make(chan request),
		standbyTimers: make(map[core.Sensor]*time.Timer),
	}
}
// RegisteredSensors returns map with sensors registered on the Device.
func (r *SensorsReader) RegisteredSensors() sensor.SensorsRegister {
	return r.sensors
}

// RegisterSensors adds given `sensors` on the SensorsReader sensors pool.
func (r *SensorsReader) RegisterSensors(sensors ...core.Sensor) {
	for i, sensor := range sensors {
		r.sensors[sensor.ID()] = sensors[i]
	}
}

// UnregisterSensors removes sensor by given `id` from the SensorsReader sensors pool.
func (r *SensorsReader) UnregisterSensors(ids ...string) {
	for _, id := range ids{
		if sensor, ok := r.sensors[id]; ok {
			if sensor.Active() {
				sensor.Close()
			}
			delete(r.sensors, id)
		}
	}
}

// SubscribeReceiver creates receiver subscription routine with given `handler`
// and starts creating sensor reading requests every given `interval`.
func (r *SensorsReader) SubscribeReceiver(
	ctx context.Context,
	handler ReceiverFunc,
	interval time.Duration,
	metrics ...models.Metric,
) context.CancelFunc {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		for {
			r.requests <- request{
				Metrics: metrics,
				Handler: handler,
			}

			select {
			case <- ctx.Done():
				return
			default:
				time.Sleep(interval)
			}
		}
	}()

	return cancel
}

// SendRequest creates single request for sensor readings that will be handled with given `handler`.
func (r *SensorsReader) SendRequest(handler ReceiverFunc, metrics ...models.Metric) {
	r.requests <- request{
		Metrics: metrics,
		Handler: handler,
	}
}

// Run starts working on the on the received requests by reading sensors data.
func (r *SensorsReader) Run(ctx context.Context) {
	ctx, r.cancel = context.WithCancel(ctx)
	r.active = true

	go r.once.Do(func() {
		for {
			select {
			case request := <- r.requests:
				go r.handleRequest(ctx, request)
			case <- ctx.Done():
				shared.Logger.Debug("Sensors reader process ended.")
				return
			}
		}
	})
}

// Active determines whether the SensorReader instance is running.
func (r *SensorsReader) Active() bool {
	return r.active
}

// Close stops SensorReader working routine and clears allocated resources.
func (r *SensorsReader) Close() {
	r.active = false
	r.cancel()

	for _, sensor := range r.sensors {
		if sensor.Active() {
			if err := sensor.Close(); err != nil {
				shared.Logger.Error(errors.Wrapf(err, "failed to close connection to '%s' sensor", sensor.ID()))
			}
		}
	}
}

func (r *SensorsReader) handleRequest(ctx context.Context, req request) {
	var (
		waitGroup = &sync.WaitGroup{}
		pipe = make(sensor.ReadingsPipe)
	)


	for _, metric := range req.Metrics {
		pipe[metric] = make(chan sensor.ReadingResult, 3)
	}

	for _, sn := range r.sensors {
		for _, metric := range req.Metrics {
			if suitable(sn, metric) {
				ctx, cancel := context.WithTimeout(ctx, 2 * time.Second) // TODO: configure or base on request period
				sensorCtx := sensor.NewReaderContext(ctx, sn)
				sensorCtx.Pipe = pipe

				if err := r.initSensor(sn); err != nil {
					sensorCtx.Error(err)
					continue
				}

				waitGroup.Add(1)
				defer cancel()
				go r.readSensor(sensorCtx, sn, waitGroup)

				break
			}
		}
	}

	waitGroup.Wait()

	results := aggregate(pipe)
	req.Handler(results)

	return
}

func suitable(sensor core.Sensor, metric models.Metric) bool {
	for _, m := range sensor.Metrics() {
		if metric == m {
			return true
		}
	}

	return false
}

func aggregate(pipe sensor.ReadingsPipe) ReadingResults {
	var (
		results = make(ReadingResults)
	)

	for metric, ch := range pipe {
		readings := make([]sensor.ReadingResult, 0)

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

func (r *SensorsReader) initSensor(sn core.Sensor) error {
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

func (r *SensorsReader) readSensor(ctx *sensor.Context, sn core.Sensor, wg *sync.WaitGroup) {
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

func handleStandby(t *time.Timer, sn core.Sensor) {
	<-t.C
	sn.Close()
}

func selectResult(results []sensor.ReadingResult) (result float64) {
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
