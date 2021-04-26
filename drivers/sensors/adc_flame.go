package sensors

import (
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripherals"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
)

type ADCFlame struct {
	peripherals.ADC
	samples int
}

func NewADCFlame(addr uint16, bus int) sensor.Sensor {
	return &ADCFlame{
		ADC: peripherals.NewADC(addr, bus, peripherals.WithConversion(func(raw float64) float64 {
			volts := raw / ADS1115_SAMPLES_PER_READ * ADS1115_VOLTS_PER_SAMPLE
			return volts
		}), peripherals.WithBias(ADC_FLAME_BIAS)),
		samples: viper.GetInt("sensors.analog.samples_per_read"),
	}
}

func (s *ADCFlame) ID() string {
	return "ADC_Flame"
}

func (s *ADCFlame) Read() float64 {
	return s.RMS(s.samples, nil)
}

func (s *ADCFlame) Harvest(ctx *sensor.Context) {
	ctx.For(metrics.Flame).Write(s.Read())
}

func (s *ADCFlame) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Flame,
	}
}
