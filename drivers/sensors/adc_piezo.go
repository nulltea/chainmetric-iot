package sensors

import (
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripherals"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
)

type ADCPiezo struct {
	peripherals.ADC
}

func NewADCPiezo(addr uint16, bus int) sensor.Sensor {
	return &ADCPiezo{
		ADC: peripherals.NewADC(addr, bus, peripherals.WithConversion(func(raw float64) float64 {
			volts := raw / ADS1115_SAMPLES_PER_READ * ADS1115_VOLTS_PER_SAMPLE
			return volts
		})),
	}
}

func (s *ADCPiezo) ID() string {
	return "ADC_Piezo"
}

func (s *ADCPiezo) Read() float64 {
	return s.RMS(100, nil)
}

func (s *ADCPiezo) Harvest(ctx *sensor.Context) {
	ctx.For(metrics.Vibration).Write(s.Read())
}

func (s *ADCPiezo) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Vibration,
	}
}
