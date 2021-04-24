package sensors

import (
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/peripherals"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
)

type ADCMQ9 struct {
	peripherals.ADC
}

func NewADCMQ9(addr uint16, bus int) *ADCMQ9 {
	return &ADCMQ9{
		ADC: peripherals.NewADC(addr, bus),
	}
}

func (s *ADCMQ9) ID() string {
	return "ADC-MQ9"
}

func (s *ADCMQ9) Harvest(ctx *Context) {
	ctx.For(metrics.AirPetroleumConcentration).WriteWithError(s.ReadRetry(5))
}

func (s *ADCMQ9) Metrics() []models.Metric {
	return []models.Metric {
		metrics.AirPetroleumConcentration,
	}
}
