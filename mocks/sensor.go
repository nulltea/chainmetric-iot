package mocks

import (
	"math/rand"
	"time"

	"sensorsys/model"
	"sensorsys/readings/sensor"
)

type MockSensor struct {
	duration time.Duration
	metrics []model.Metric
}

func NewMockSensor(duration time.Duration, metrics ...model.Metric) *MockSensor {
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

func (s *MockSensor) Metrics() []model.Metric {
	return s.metrics
}

func (s *MockSensor) Active() bool {
	return true
}

func (s *MockSensor) Close() error {
	return nil
}
