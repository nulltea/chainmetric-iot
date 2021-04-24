package sensors

import (
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/peripherals"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/sensor"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
)

type ADCFlame struct {
	peripherals.ADC
}

func NewADCFlame(addr uint16, bus int) sensor.Sensor {
	return &ADCFlame{
		ADC: peripherals.NewADC(addr, bus),
	}
}

func (s *ADCFlame) ID() string {
	return "ADC_Flame"
}

func (s *ADCFlame) Harvest(ctx *sensor.Context) {
	ctx.For(metrics.Flame).WriteWithError(s.ReadRetry(5))
}

func (s *ADCFlame) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Flame,
	}
}
