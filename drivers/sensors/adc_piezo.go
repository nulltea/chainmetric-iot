package sensors

import (
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripherals"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
)

type ADCPiezo struct {
	peripherals.ADC
	samples int
}

func NewADCPiezo(addr uint16, bus int) sensor.Sensor {
	return &ADCPiezo{
		ADC: peripherals.NewADC(addr, bus, peripherals.WithConversion(func(raw float64) float64 {
			volts := raw / peripherals.ADS1115_SAMPLES_PER_READ * peripherals.ADS1115_VOLTS_PER_SAMPLE
			return volts
		})),
		samples: viper.GetInt("sensors.analog.samples_per_read"),
	}
}

func (s *ADCPiezo) ID() string {
	return "ADC_Piezo"
}

func (s *ADCPiezo) Read() float64 {
	return s.RMS(s.samples, nil)
}

func (s *ADCPiezo) Harvest(ctx *sensor.Context) {
	ctx.For(metrics.Vibration).Write(s.Read())
}

func (s *ADCPiezo) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Vibration,
	}
}
