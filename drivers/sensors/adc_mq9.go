package sensors

import (
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripherals"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
)

type ADCMQ9 struct {
	peripherals.ADC
}

func NewADCMQ9(addr uint16, bus int) sensor.Sensor {
	return &ADCMQ9{
		ADC: peripherals.NewADC(addr, bus),
	}
}

func (s *ADCMQ9) ID() string {
	return "ADC-MQ9"
}

func (s *ADCMQ9) Harvest(ctx *sensor.Context) {
	ctx.For(metrics.AirPetroleumConcentration).WriteWithError(s.ReadRetry(5))
}

func (s *ADCMQ9) Metrics() []models.Metric {
	return []models.Metric {
		metrics.AirPetroleumConcentration,
	}
}
