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
		ADC: peripherals.NewADC(addr, bus, peripherals.WithConversion(func(raw float64) float64 {
			volts := raw / ADS1115_SAMPLES_PER_READ * ADS1115_VOLTS_PER_SAMPLE
			return volts * 1000 / ADC_HALL_SENSITIVITY
		}), peripherals.WithBias(ADC_HALL_BIAS)),
	}
}

func (s *ADCHall) ID() string {
	return "ADC_Hall"
}

func (s *ADCHall) Read() float64 {
	return s.RMS(100, nil)
}

func (s *ADCHall) Harvest(ctx *sensor.Context) {
	ctx.For(metrics.Magnetism).Write(s.Read())
}

func (s *ADCHall) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Magnetism,
	}
}
