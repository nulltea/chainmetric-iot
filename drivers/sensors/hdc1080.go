package sensors

import (
	"fmt"
	"sync"
	"time"

	"github.com/d2r2/go-i2c"
	"github.com/pkg/errors"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
)

type HDC1080 struct {
	addr     uint8
	bus      int
	i2c      *i2c.I2C
	attempts int
}

func NewHDC1080(addr uint8, bus int) *HDC1080 {
	return &HDC1080{
		addr:     addr,
		bus:      bus,
		attempts: 10,
	}
}

func (s *HDC1080) ID() string {
	return "HDC1080"
}

func (s *HDC1080) Init() (err error) {
	s.i2c, err = i2c.NewI2C(s.addr, s.bus); if err != nil {
		return
	}

	if !s.Verify() {
		return fmt.Errorf("not HDC1080 sensor")
	}

	_, err = s.i2c.WriteBytes([]byte{HDC1080_CONFIGURATION_REGISTER, HDC1080_CONFIG_ACQUISITION_MODE >> 8, 0x00})
	time.Sleep(15 * time.Millisecond)

	return
}

func (s *HDC1080) ReadTemperature() (float64, error) {
	if _, err := s.i2c.WriteBytes([]byte{HDC1080_TEMPERATURE_REGISTER}); err != nil {
		return 0, errors.Wrap(err, "failed write to temperature register")
	}

	var (
		data = make([]byte, 2)
		left = s.attempts
		err error
	)

	for left >= 0 {
		left--
		time.Sleep(65 * time.Millisecond)

		if _, err = s.i2c.ReadBytes(data); err != nil {
			continue
		}

		raw := float64(int(data[0]) << 8 + int(data[1]))

		return (raw / 65536.0) * 165.0 - 40.0, nil
	}

	return 0, errors.Wrap(err, "failed read from temperature register")
}

func (s *HDC1080) ReadHumidity() (float64, error) {
	if _, err := s.i2c.WriteBytes([]byte{HDC1080_HUMIDITY_REGISTER}); err != nil {
		return 0, errors.Wrap(err, "failed write to humidity register")
	}

	var (
		data = make([]byte, 2)
		left = s.attempts
		err error
	)

	for left >= 0 {
		left--
		time.Sleep(65 * time.Millisecond)

		if _, err = s.i2c.ReadBytes(data); err != nil {
			continue
		}

		raw := float64(int(data[0]) << 8 + int(data[1]))

		return (raw / 65536.0) * 100.0, nil
	}

	return 0, errors.Wrap(err, "failed read from humidity register")
}

func (s *HDC1080) Harvest(ctx *Context) {
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

func (s *HDC1080) Active() bool {
	return s.i2c != nil
}

func (s *HDC1080) Close() error {
	defer s.clean()
	return s.i2c.Close()
}

func (s *HDC1080) clean() {
	s.i2c = nil
}
