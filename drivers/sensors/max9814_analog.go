package sensors

import (
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/peripherals"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
)

type AnalogMAX9814 struct {
	ch *peripherals.AnalogChannel
	active bool
}

func NewAnalogMAX9814(ch *peripherals.AnalogChannel) *AnalogMAX9814 {
	return &AnalogMAX9814{
		ch: ch,
	}
}

func (s *AnalogMAX9814) ID() string {
	return "Analog_MAX9814"
}

func (s *AnalogMAX9814) Init() error {
	s.active = true
	return nil
}

func (s *AnalogMAX9814) Harvest(ctx *Context) {
	ctx.For(metrics.NoiseLevel).Write(s.ch.Get())
}

func (s *AnalogMAX9814) Metrics() []models.Metric {
	return []models.Metric {
		metrics.NoiseLevel,
	}
}

func (s *AnalogMAX9814) Active() bool {
	return s.active
}

func (s *AnalogMAX9814) Close() error {
	s.active = false
	return nil
}
