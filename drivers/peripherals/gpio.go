package peripherals

import (
	"fmt"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"

	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// GPIO provides wrapper for GPIO peripheral
type GPIO struct {
	gpio.PinIO
	pin int
}

func NewGPIO(pin int) *GPIO {
	return &GPIO{
		pin: pin,
	}
}

func (g *GPIO) Init() error {
	var (
		name = shared.NtoPinName(g.pin)
		pin = gpioreg.ByName(name)
	)

	if pin == gpio.INVALID || pin == nil {
		return fmt.Errorf("pin %s is invalid", name)
	}

	g.PinIO = pin

	return nil
}

func (g *GPIO) High() error {
	return g.Out(gpio.High)
}

func (g *GPIO) Low() error {
	return g.Out(gpio.Low)
}

func (g *GPIO) IsHigh() bool {
	return g.Read() == gpio.High
}

func (g *GPIO) IsLow() bool {
	return g.Read() == gpio.Low
}
