package peripherals

import (
	"github.com/pkg/errors"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"

	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

// I2C provides wrapper for I2C peripheral
type I2C struct {
	i2c.Dev
	name string
	bus  i2c.BusCloser
}

func NewI2C(addr uint16, bn int) *I2C {
	return &I2C{
		Dev: i2c.Dev{
			Addr: addr,
		},
		name: shared.NtoI2cBusName(bn),
	}
}

func (i *I2C) Init() (err error) {
	if i.bus, err = i2creg.Open(i.name); err != nil {
		return errors.Wrapf(err, "failed to open an I2C bus on %s", i.name)
	}

	i.Bus = i.bus

	return
}
