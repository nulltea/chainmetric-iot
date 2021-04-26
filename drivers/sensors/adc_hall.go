package sensors

import (
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripherals"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
)

type ADCHall struct {
	peripherals.ADC
}

func NewADCHall(addr uint16, bus int) sensor.Sensor {
	return &ADCHall{
		ADC: peripherals.NewADC(addr, bus),
	}
}

func (s *ADCHall) ID() string {
	return "ADC_Hall"
}

func (s *ADCHall) Read() float64 {
	return s.Aggregate(100, nil) - ADC_HALL_BIAS
}

func (s *ADCHall) Harvest(ctx *sensor.Context) {
	ctx.For(metrics.Magnetism).Write(s.Read())
}

func (s *ADCHall) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Magnetism,
	}
}
