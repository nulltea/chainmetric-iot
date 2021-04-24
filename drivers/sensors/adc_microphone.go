package sensors

import (
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/peripherals"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/sensor"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
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

func (s *ADCMicrophone) Harvest(ctx *sensor.Context) {
	ctx.For(metrics.NoiseLevel).WriteWithError(s.ReadRetry(5))
}

func (s *ADCMicrophone) Metrics() []models.Metric {
	return []models.Metric {
		metrics.NoiseLevel,
	}
}
