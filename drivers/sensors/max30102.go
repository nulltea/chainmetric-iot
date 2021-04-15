package sensors

import (
	"github.com/cgxeiji/max3010x"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

type MAX30102 struct {
	addr uint8
	bus int
	dev *max3010x.Device
}

func NewMAX30102(addr uint8, bus int) *MAX30102 {
	return &MAX30102{
		addr: addr,
		bus: bus,
	}
}

func (s *MAX30102) ID() string {
	return "MAX30102"
}

func (s *MAX30102) Init() (err error) {
	s.dev, err = max3010x.New(
		max3010x.WithSpecificBus(shared.NtoI2cBusName(s.bus)),
	); if err != nil {
		return
	}

	if err = s.dev.Startup(); err != nil {
		return err
	}

	return
}

func (s *MAX30102) Harvest(ctx *Context) {
	ctx.For(metrics.HeartRate).WriteWithError(s.ReadHeartRate())
	ctx.For(metrics.BloodOxidation).WriteWithError(s.ReadSpO2())
}

func (s *MAX30102) Metrics() []models.Metric {
	return []models.Metric {
		metrics.HeartRate,
		metrics.BloodOxidation,
	}
}

func (s *MAX30102) ReadHeartRate() (float64, error) {
	return s.dev.HeartRate()
}

func (s *MAX30102) ReadSpO2() (float64, error) {
	return s.dev.SpO2()
}

func (s *MAX30102) Active() bool {
	return s.dev != nil
}

// Close disconnects from the device
func (s *MAX30102) Close() error {
	s.dev.Close()
	s.dev = nil
	return nil
}
