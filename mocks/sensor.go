package mocks

import (
	"math/rand"
	"time"

	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/engine/sensor"
)

type MockSensor struct {
	duration time.Duration
	metrics []models.Metric
}

func NewMockSensor(duration time.Duration, metrics ...models.Metric) *MockSensor {
	return &MockSensor{
		duration,
		metrics,
	}
}

func (s *MockSensor) ID() string {
	return "mock-sensor"
}

func (s *MockSensor) Init() error {
	rand.Seed(time.Now().UnixNano())
	return nil
}

func (s *MockSensor) Harvest(ctx *sensor.Context) {
	time.Sleep(s.duration)

	for _, metric := range s.metrics {
		ctx.For(metric).Write(rand.Float32())
	}
}

func (s *MockSensor) Metrics() []models.Metric {
	return s.metrics
}

func (s *MockSensor) Active() bool {
	return true
}

func (s *MockSensor) Close() error {
	return nil
}
