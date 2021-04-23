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

func (s *AnalogPZT) Read() float64 {
	var (
		v uint64
		i int
	)

	for i != 100 {
		if vc := s.ch.Get(); vc != 0 {
			v += uint64(vc)
			i++
		}
	}

	return float64(v / 100)
}

func (s *AnalogPZT) Harvest(ctx *Context) {
	ctx.For(metrics.Vibration).Write(s.Read())
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
