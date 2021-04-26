package sensors

import (
	"math"

	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripherals"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
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

func (s *ADCFlame) Read() float64 {
	return math.Abs(s.Aggregate(100, nil) - ADC_FLAME_BIAS)
}

func (s *ADCFlame) Harvest(ctx *sensor.Context) {
	ctx.For(metrics.Flame).Write(s.Read())
}

func (s *ADCFlame) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Flame,
	}
}
