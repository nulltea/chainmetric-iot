package sensors

import (
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/peripherals"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
)

type ADCPiezo struct {
	peripherals.ADC
}

func NewADCPiezo(addr uint16, bus int) *ADCPiezo {
	return &ADCPiezo{
		ADC: peripherals.NewADC(addr, bus),
	}
}

func (s *ADCPiezo) ID() string {
	return "ADC_Piezo"
}

func (s *ADCPiezo) Harvest(ctx *Context) {
	ctx.For(metrics.Vibration).WriteWithError(s.ReadRetry(5))
}

func (s *ADCPiezo) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Vibration,
	}
}