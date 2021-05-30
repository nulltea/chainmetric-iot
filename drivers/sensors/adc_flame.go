package sensors

import (
	"sync"

	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-sensorsys/core/dev/sensor"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/periphery"
)

var (
	adcFlameMutex = &sync.Mutex{}
)

type ADCFlame struct {
	periphery.ADC
	samples int
}

func NewADCFlame(addr uint16, bus int) sensor.Sensor {
	return &ADCFlame{
		ADC: periphery.NewADC(addr, bus, periphery.WithConversion(func(raw float64) float64 {
			volts := raw / periphery.ADS1115_SAMPLES_PER_READ * periphery.ADS1115_VOLTS_PER_SAMPLE
			return volts
		}), periphery.WithBias(ADC_FLAME_BIAS), periphery.WithI2CMutex(adcFlameMutex)),
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
	ctx.WriterFor(metrics.Flame).Write(s.Read())
}

func (s *ADCFlame) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Flame,
	}
}
