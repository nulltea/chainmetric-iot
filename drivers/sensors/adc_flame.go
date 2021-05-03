package sensors

import (
	"sync"

	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripheries"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
)

var (
	adcFlameMutex = sync.Mutex{}
)

type ADCFlame struct {
	peripheries.ADC
	samples int
}

func NewADCFlame(addr uint16, bus int) sensor.Sensor {
	return &ADCFlame{
		ADC: peripheries.NewADC(addr, bus, peripheries.WithConversion(func(raw float64) float64 {
			volts := raw / peripheries.ADS1115_SAMPLES_PER_READ * peripheries.ADS1115_VOLTS_PER_SAMPLE
			return volts
		}), peripheries.WithBias(ADC_FLAME_BIAS)),
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
