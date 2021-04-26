package sensors

import (
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripherals"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
)

type ADCMicrophone struct {
	peripherals.ADC
}

func NewADCMicrophone(addr uint16, bus int) sensor.Sensor {
	return &ADCMicrophone{
		ADC: peripherals.NewADC(addr, bus),
	}
}

func (s *ADCMicrophone) ID() string {
	return "ADC_Microphone"
}

func (s *ADCMicrophone) Read() float64 {
	return ADC_MICROPHONE_REGRESSION_C1* (s.Aggregate(100, nil) -ADC_MICROPHONE_BIAS) +
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
