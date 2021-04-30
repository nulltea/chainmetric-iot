package sensors

import (
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripheries"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
)

type ADCMicrophone struct {
	peripheries.ADC
	samples int
}

func NewADCMicrophone(addr uint16, bus int) sensor.Sensor {
	return &ADCMicrophone{
		ADC:     peripheries.NewADC(addr, bus),
		samples: viper.GetInt("sensors.analog.samples_per_read"),
	}
}

func (s *ADCMicrophone) ID() string {
	return "ADC_Microphone"
}

func (s *ADCMicrophone) Read() float64 {
	return ADC_MICROPHONE_REGRESSION_C1 * (s.RMS(s.samples, nil) - ADC_MICROPHONE_BIAS) +
		ADC_MICROPHONE_REGRESSION_C2
}

func (s *ADCMicrophone) Harvest(ctx *sensor.Context) {
	ctx.For(metrics.NoiseLevel).Write(s.Read())
}

func (s *ADCMicrophone) Metrics() []models.Metric {
	return []models.Metric {
		metrics.NoiseLevel,
	}
}
