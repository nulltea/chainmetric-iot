package sensors

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripherals"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
)

type HDC1080 struct {
	*peripherals.I2C
	attempts int
}

func NewHDC1080(addr uint16, bus int) sensor.Sensor {
	return &HDC1080{
		I2C:      peripherals.NewI2C(addr, bus),
		attempts: 10,
	}
}

func (s *HDC1080) ID() string {
	return "HDC1080"
}

func (s *HDC1080) Init() error {
	if err := s.I2C.Init(); err != nil {
		return err
	}

	if !s.Verify() {
		return fmt.Errorf("driver is not compatiple with specified sensor")
	}

	if err := s.WriteRegBytes(HDC1080_CONFIGURATION_REGISTER, HDC1080_CONFIG_ACQUISITION_MODE >> 8, 0x00); err != nil {
		return err
	}

	time.Sleep(15 * time.Millisecond)

	return nil
}

func (s *HDC1080) ReadTemperature() (float64, error) {
	if err := s.WriteBytes(HDC1080_TEMPERATURE_REGISTER); err != nil {
		return 0, errors.Wrap(err, "failed write to temperature register")
	}

	var (
		data []byte
		left = s.attempts
		err error
	)

	for left >= 0 {
		left--
		time.Sleep(65 * time.Millisecond)

		if data, err = s.ReadBytes(2); err != nil {
			continue
		}

		raw := float64(int(data[0]) << 8 + int(data[1]))

		return (raw / 65536.0) * 165.0 - 40.0, nil
	}

	return 0, errors.Wrap(err, "failed read from temperature register")
}

func (s *HDC1080) ReadHumidity() (float64, error) {
	if err := s.WriteBytes(HDC1080_HUMIDITY_REGISTER); err != nil {
		return 0, errors.Wrap(err, "failed write to humidity register")
	}

	var (
		data []byte
		left = s.attempts
		err error
	)

	for left >= 0 {
		left--
		time.Sleep(65 * time.Millisecond)

		if data, err = s.ReadBytes(2); err != nil {
			continue
		}

		raw := float64(int(data[0]) << 8 + int(data[1]))

		return (raw / 65536.0) * 100.0, nil
	}

	return 0, errors.Wrap(err, "failed read from humidity register")
}

func (s *HDC1080) Harvest(ctx *sensor.Context) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		ctx.For(metrics.Temperature).WriteWithError(s.ReadTemperature())
		wg.Done()
	}()

	go func() {
		ctx.For(metrics.Humidity).WriteWithError(s.ReadHumidity())
		wg.Done()
	}()

	wg.Wait()
}

func (s *HDC1080) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Temperature,
		metrics.Humidity,
	}
}

func (s *HDC1080) Verify() bool {
	return true
}
