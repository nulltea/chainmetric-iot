package engine

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-iot/core/dev/sensor"
	"github.com/timoth-y/chainmetric-iot/shared"
)

type (
	// SensorsReader defines operational structure of the sensors reading engine.
	SensorsReader struct {
		once          *sync.Once
		sensors       sensor.SensorsRegister
		requests      chan request
		standbyTimers map[sensor.Sensor]*time.Timer
		active        bool
		cancel        context.CancelFunc
	}

	// ReadingResults defines map of values collected from sensor.Sensor for requested models.Metrics.
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
		sensors:       make(map[string]sensor.Sensor),
		requests:      make(chan request),
		standbyTimers: make(map[sensor.Sensor]*time.Timer),
	}
}
// RegisteredSensors returns map with sensors registered on the engine.SensorsReader.
func (r *SensorsReader) RegisteredSensors() sensor.SensorsRegister {
	return r.sensors
}

// RegisterSensors adds given `sensors` on the SensorsReader sensors pool.
func (r *SensorsReader) RegisterSensors(sensors ...sensor.Sensor) {
	for i, s := range sensors {
		r.sensors[s.ID()] = sensors[i]
	}
}

// UnregisterSensors removes sensor by given `id` from the SensorsReader sensors pool.
func (r *SensorsReader) UnregisterSensors(ids ...string) {
	for _, id := range ids{
		if s, ok := r.sensors[id]; ok {
			if s.Active() {
				if err := s.Close(); err != nil {
					shared.Logger.Error(errors.Wrapf(err, "failed to close connection to '%s' sensor", s.ID()))
				}
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
	go func(ctx context.Context) {
		LOOP: for {
			r.requests <- request{
				Metrics: metrics,
				Handler: handler,
			}

			select {
			case <- ctx.Done():
				break LOOP
			default:
				time.Sleep(interval)
			}
		}
	}(ctx)

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
				shared.Logger.Debug("Sensors reader engine routine ended")
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

	for _, s := range r.sensors {
		if s.Active() {
			if err := s.Close(); err != nil {
				shared.Logger.Error(errors.Wrapf(err, "failed to close connection to '%s' sensor", s.ID()))
			}
		}
	}
}

func (r *SensorsReader) handleRequest(ctx context.Context, req request) {
	var (
		waitGroup = &sync.WaitGroup{}
		pipe = make(sensor.ReadingsPipe)
	)

	// Init channels in request results pipe:
	for _, metric := range req.Metrics {
		pipe[metric] = make(chan sensor.ReadingResult, 3)
	}

	// Create single timeout context for all suitable for requested metrics sensors:
	// TODO: each sensor may take different amount of time to be read,
	//  thus some kind of deterministic timeout handling approach is required here.
	//  Readings interval specified by receiver also should be taken in the account here.
	ctx, cancel := context.WithTimeout(ctx, 3 * time.Second)
	defer cancel()

	// Go through available sensors to check is there any compatible ones for requested metrics,
	// and if so perform reading from them:
	for _, sn := range r.sensors {
		for _, metric := range req.Metrics {
			if suitable(sn, metric) {
				waitGroup.Add(1)

				go func(sn sensor.Sensor) {
					// Create new reading context for sensor and assign channels pipe,
					// where reading results will be dumped into:
					sensorCtx := sensor.NewReaderContext(ctx, sn)
					sensorCtx.Pipe = pipe

					// First time use initialization along with stand by handling:
					if err := r.initSensor(sn); err != nil {
						sensorCtx.Error(err)
						return
					}

					r.readSensor(sensorCtx, sn, waitGroup)
				}(sn)

				break
			}
		}
	}

	// Wait until all required sensors finish being read or timed out:
	waitGroup.Wait()

	// Finally, aggregate sensor reading results and handle them by passing to receiver:
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
	shared.Execute(sn.Close, fmt.Sprintf("failed to close connection to '%s' sensor", sn.ID()))
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
