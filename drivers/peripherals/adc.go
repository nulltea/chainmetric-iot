package peripherals

import (
	"github.com/MichaelS11/go-ads"
	"github.com/pkg/errors"

	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

// ADC defines analog to digital peripheral interface.
type ADC interface {
	Init() error
	Read() (uint16, error)
	ReadRetry(int) (uint16, error)
	Active() bool
	Close() error
}

// ADS1115 implements ADC driver for ADS1115 device.
type ADS1115 struct {
	*ads.ADS
	Addr   uint16
	Bus    string
	active bool
}

// NewADC returns a new ADC implementation via ADS1115 device driver.
func NewADC(addr uint16, bus int) *ADS1115 {
	return &ADS1115{
		Bus: shared.NtoI2cBusName(bus),
		Addr: addr,
	}
}

// Init sets up the device for communication.
func (d *ADS1115) Init() (err error) {
	if d.ADS, err = ads.NewADS(d.Bus, d.Addr, "ADS1115"); err != nil {
		return errors.Wrapf(err, "failed to init ADS1115 device on '%s' bus and 0x%X address", d.Bus, d.Addr)
	}

	d.active = true

	return nil
}

func (d *ADS1115) Active() bool {
	return d.active
}

func (d *ADS1115) Close() error {
	d.active = false
	return d.ADS.Close()
}
