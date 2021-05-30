package periphery

import (
	"github.com/pkg/errors"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"

	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// GPIO provides wrapper for GPIO peripheral.
type GPIO struct {
	gpio.PinIO
	pin string
}

// NewGPIO constructs new GPIO driver instance.
func NewGPIO(pin int) *GPIO {
	return &GPIO{
		pin: shared.NtoPinName(pin),
	}
}

// Init performs GPIO driver initialization.
func (g *GPIO) Init() error {
	var (
		pin = gpioreg.ByName(g.pin)
	)

	if pin == gpio.INVALID || pin == nil {
		return errors.Errorf("pin %s is invalid", g.pin)
	}

	g.PinIO = pin

	if err := g.Low(); err != nil {
		return errors.Wrapf(err, "failed initialising pin", g.pin)
	}

	return nil
}

// High sends high level signal to GPIO pin.
func (g *GPIO) High() error {
	return g.Out(gpio.High)
}

// Low sends low level signal to GPIO pin.
func (g *GPIO) Low() error {
	return g.Out(gpio.Low)
}

// IsHigh determines whether the GPIO pin is on High state.
func (g *GPIO) IsHigh() bool {
	return g.Read() == gpio.High
}

// IsLow determines whether the GPIO pin is on Low state.
func (g *GPIO) IsLow() bool {
	return g.Read() == gpio.Low
}
