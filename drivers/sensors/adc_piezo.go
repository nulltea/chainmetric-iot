package sensors

import (
	"sync"

	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripheries"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

var (
	adcPiezoMutex = sync.Mutex{}
)

type ADCPiezo struct {
	peripheries.ADC
	samples int
}

func NewADCPiezo(addr uint16, bus int) sensor.Sensor {
	return &ADCPiezo{
		ADC: peripheries.NewADC(addr, bus, peripheries.WithConversion(func(raw float64) float64 {
			shared.Logger.Debug("ADC_Piezo", "-> raw =", raw)
			volts := raw / peripheries.ADS1115_SAMPLES_PER_READ * peripheries.ADS1115_VOLTS_PER_SAMPLE
			shared.Logger.Debug("ADC_Piezo", "-> volts =", volts)
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
