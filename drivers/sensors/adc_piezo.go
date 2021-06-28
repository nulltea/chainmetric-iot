package sensors

import (
	"sync"

	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-iot/core/dev/sensor"
	"github.com/timoth-y/chainmetric-iot/drivers/periphery"
	"github.com/timoth-y/chainmetric-iot/shared"
)

var (
	adcPiezoMutex = &sync.Mutex{}
)

type ADCPiezo struct {
	periphery.ADC
	samples int
}

func NewADCPiezo(addr uint16, bus int) sensor.Sensor {
	return &ADCPiezo{
		ADC: periphery.NewADC(addr, bus, periphery.WithConversion(func(raw float64) float64 {
			shared.Logger.Debug("ADC_Piezo", "-> raw =", raw)
			volts := raw / periphery.ADS1115_SAMPLES_PER_READ * periphery.ADS1115_VOLTS_PER_SAMPLE
			shared.Logger.Debug("ADC_Piezo", "-> volts =", volts)
			return volts
		}), periphery.WithI2CMutex(adcPiezoMutex)),
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
	ctx.WriterFor(metrics.Vibration).Write(s.Read())
}

func (s *ADCPiezo) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Vibration,
	}
}
