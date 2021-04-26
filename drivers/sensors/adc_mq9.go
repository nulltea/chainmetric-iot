package sensors

import (
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripherals"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
)

type ADCMQ9 struct {
	peripherals.ADC
}

func NewADCMQ9(addr uint16, bus int) sensor.Sensor {
	return &ADCMQ9{
		ADC: peripherals.NewADC(addr, bus, peripherals.WithConversion(func(raw float64) float64 {
			volts := raw / ADS1115_SAMPLES_PER_READ * ADS1115_VOLTS_PER_SAMPLE
			resAir := (ADC_MQ9_RESISTANCE - volts) / volts
			return resAir / ADC_MQ9_SENSITIVITY * 1000
		}), peripherals.WithBias(ADC_MQ9_BIAS)),
	}
}

func (s *ADCMQ9) ID() string {
	return "ADC-MQ9"
}

func (s *ADCMQ9) Read() float64 {
	return s.RMS(100, nil)
}

func (s *ADCMQ9) Harvest(ctx *sensor.Context) {
	ctx.For(metrics.AirPetroleumConcentration).Write(s.Read())
}

func (s *ADCMQ9) Metrics() []models.Metric {
	return []models.Metric {
		metrics.AirPetroleumConcentration,
	}
}
