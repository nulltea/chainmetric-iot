package sensors

import (
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/peripherals"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/sensor"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
)

type ADCHall struct {
	peripherals.ADC
}

func NewADCHall(addr uint16, bus int) sensor.Sensor {
	return &ADCHall{
		ADC: peripherals.NewADC(addr, bus),
	}
}

func (s *ADCHall) ID() string {
	return "ADC_Hall"
}

func (s *ADCHall) Harvest(ctx *sensor.Context) {
	ctx.For(metrics.Magnetism).WriteWithError(s.ReadRetry(5))
}

func (s *ADCHall) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Magnetism,
	}
}
