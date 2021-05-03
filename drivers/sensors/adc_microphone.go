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
	adcMicMutex = sync.Mutex{}
)

type ADCMic struct {
	peripheries.ADC
	samples int
}

func NewADCMicrophone(addr uint16, bus int) sensor.Sensor {
	return &ADCMic{
		ADC:     peripheries.NewADC(addr, bus),
		samples: viper.GetInt("sensors.analog.samples_per_read"),
	}
}

func (s *ADCMic) ID() string {
	return "ADC_Microphone"
}

func (s *ADCMic) Read() float64 {
	return ADC_MICROPHONE_REGRESSION_C1 * (s.RMS(s.samples, nil) - ADC_MICROPHONE_BIAS) +
		ADC_MICROPHONE_REGRESSION_C2
}

func (s *ADCMic) Harvest(ctx *sensor.Context) {
	ctx.For(metrics.NoiseLevel).Write(s.Read())
}

func (s *ADCMic) Metrics() []models.Metric {
	return []models.Metric {
		metrics.NoiseLevel,
	}
}
