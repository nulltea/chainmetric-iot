package sensors

import (
	"sync"

	"github.com/cgxeiji/max3010x"
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripheries"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

var (
	max30102Mutex = &sync.Mutex{}
)

type MAX30102 struct {
	*max3010x.Device
	i2c *peripheries.I2C
	addr uint16
	bus int
}

func NewMAX30102(addr uint16, bus int) sensor.Sensor {
	return &MAX30102{
		i2c:  peripheries.NewI2C(addr, bus, peripheries.WithMutex(max30102Mutex)),
		addr: addr,
		bus:  bus,
	}
}

func (s *MAX30102) ID() string {
	return "MAX30102"
}

func (s *MAX30102) Init() (err error) {
	s.Device, err = max3010x.New(
		max3010x.WithSpecificBus(shared.NtoI2cBusName(s.bus)),
		max3010x.WithAddress(s.addr),
	); if err != nil {
		return
	}

	if err = s.Startup(); err != nil {
		return err
	}

	return
}

func (s *MAX30102) Harvest(ctx *sensor.Context) {
	ctx.For(metrics.HeartRate).WriteWithError(s.HeartRate())
	ctx.For(metrics.BloodOxidation).WriteWithError(s.SpO2())
}

func (s *MAX30102) Metrics() []models.Metric {
	return []models.Metric {
		metrics.HeartRate,
		metrics.BloodOxidation,
	}
}

func (s *MAX30102) Verify() bool {
	if !s.i2c.Verify() {
		return false
	}

	if devID, err := s.i2c.ReadReg(MAX30102_DEVICE_ID_REGISTER); err == nil {
		return devID == MAX30102_DEVICE_ID
	}

	return false
}

func (s *MAX30102) Active() bool {
	return s.Device != nil
}

// Close disconnects from the device
func (s *MAX30102) Close() error {
	s.Device.Close()
	s.Device = nil
	return nil
}
