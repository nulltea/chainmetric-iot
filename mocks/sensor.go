package mocks

import (
	"math/rand"
	"time"

	"sensorsys/model"
	"sensorsys/model/metrics"
	"sensorsys/readings/sensor"
)

type MockSensor struct {

}

func NewMockSensor() *MockSensor {
	return &MockSensor{}
}

func (s *MockSensor) ID() string {
	return "mock-sensor"
}

func (s *MockSensor) Init() error {
	rand.Seed(time.Now().UnixNano())
	return nil
}

func (s *MockSensor) Harvest(ctx *sensor.Context) {
	time.Sleep(100 * time.Millisecond)
	temperature, humidity, lux := rand.Float32(), rand.Float32(), rand.Int()

	ctx.For(metrics.Temperature).Write(temperature)
	ctx.For(metrics.Humidity).Write(humidity)
	ctx.For(metrics.Luminosity).Write(lux)
}

func (s *MockSensor) Metrics() []model.Metric {
	return []model.Metric {
		metrics.Temperature,
		metrics.Humidity,
		metrics.Luminosity,
	}
}

func (s *MockSensor) Close() error {
	return nil
}
