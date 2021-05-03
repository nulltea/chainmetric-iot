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
	adcMQ9Mutex = &sync.Mutex{}
)

type ADCMQ9 struct {
	peripheries.ADC
	samples int
}

func NewADCMQ9(addr uint16, bus int) sensor.Sensor {
	return &ADCMQ9{
		ADC: peripheries.NewADC(addr, bus, peripheries.WithConversion(func(raw float64) float64 {
			volts := raw / peripheries.ADS1115_SAMPLES_PER_READ * peripheries.ADS1115_VOLTS_PER_SAMPLE
			resAir := (ADC_MQ9_RESISTANCE - volts) / volts
			return resAir / ADC_MQ9_SENSITIVITY * 1000
		}), peripheries.WithBias(ADC_MQ9_BIAS), peripheries.WithI2CMutex(adcMQ9Mutex)),
		samples: viper.GetInt("sensors.analog.samples_per_read"),
	}
}

func (s *ADCMQ9) ID() string {
	return "ADC-MQ9"
}

func (s *ADCMQ9) Read() float64 {
	return s.RMS(s.samples, nil)
}

func (s *ADCMQ9) Harvest(ctx *sensor.Context) {
	ctx.For(metrics.AirPetroleumConcentration).Write(s.Read())
}

func (s *ADCMQ9) Metrics() []models.Metric {
	return []models.Metric {
		metrics.AirPetroleumConcentration,
	}
}
