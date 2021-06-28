package sensors

import (
	"math/rand"
	"time"

	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-iot/core/dev/sensor"
	"github.com/timoth-y/chainmetric-iot/drivers/periphery"
)

type (
	I2CSensorMock struct {
		*periphery.I2C
		duration time.Duration
		metrics  []models.Metric
		active   bool
	}

	StaticSensorMock struct {
		duration time.Duration
		metrics  []models.Metric
	}
)

func NewI2CSensorMock(addr uint16, bus int) sensor.Sensor {
	return &I2CSensorMock{
		I2C:      periphery.NewI2C(addr, bus),
		duration: viper.GetDuration("mocks.sensor_duration"),
		metrics:  []models.Metric{metrics.AirCO2Concentration, metrics.Luminosity, metrics.Magnetism},
	}
}

func (s *I2CSensorMock) ID() string {
	return "MOCK-I2C"
}

func (s *I2CSensorMock) Init() error {
	rand.Seed(time.Now().UnixNano())
	s.active = true
	return nil
}

func (s *I2CSensorMock) Harvest(ctx *sensor.Context) {
	time.Sleep(s.duration)

	for _, metric := range s.metrics {
		ctx.WriterFor(metric).Write(rand.Float64())
	}
}

func (s *I2CSensorMock) Metrics() []models.Metric {
	return s.metrics
}

func (s *I2CSensorMock) Verify() bool {
	return true
}

func (s *I2CSensorMock) Active() bool {
	return s.active
}

func (s *I2CSensorMock) Close() error {
	s.active = false
	return nil
}

func NewStaticSensorMock() sensor.Sensor {
	return &StaticSensorMock{
		duration: viper.GetDuration("mocks.sensor_duration"),
		metrics:  []models.Metric{metrics.Humidity, metrics.NoiseLevel, metrics.Vibration},
	}
}

func (s *StaticSensorMock) ID() string {
	return "MOCK_Static"
}

func (s *StaticSensorMock) Init() error {
	rand.Seed(time.Now().UnixNano())
	return nil
}

func (s *StaticSensorMock) Harvest(ctx *sensor.Context) {
	time.Sleep(s.duration)

	for _, metric := range s.metrics {
		ctx.WriterFor(metric).Write(rand.Float64())
	}
}

func (s *StaticSensorMock) Metrics() []models.Metric {
	return s.metrics
}

func (s *StaticSensorMock) Verify() bool {
	return true
}

func (s *StaticSensorMock) Active() bool {
	return true
}

func (s *StaticSensorMock) Close() error {
	return nil
}
