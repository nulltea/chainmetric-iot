package sensors

import (
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/peripherals"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
)

type AnalogPZT struct {
	ch *peripherals.AnalogChannel
	active bool
}

func NewAnalogPZT(ch *peripherals.AnalogChannel) *AnalogPZT {
	return &AnalogPZT{
		ch: ch,
	}
}

func (s *AnalogPZT) ID() string {
	return "Analog_PZT"
}

func (s *AnalogPZT) Init() error {
	s.active = true
	return nil
}

func (s *AnalogPZT) Harvest(ctx *Context) {
	ctx.For(metrics.Vibration).Write(s.ch.Get())
}

func (s *AnalogPZT) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Vibration,
	}
}

func (s *AnalogPZT) Active() bool {
	return s.active
}

func (s *AnalogPZT) Close() error {
	s.active = false
	return nil
}
