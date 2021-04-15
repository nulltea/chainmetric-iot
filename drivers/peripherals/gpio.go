package peripherals

import (
	"fmt"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"

	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

// GPIO provides wrapper for GPIO peripheral
type GPIO struct {
	gpio.PinIO
}

func NewGPIO(pin int) *GPIO {
	return &GPIO{
		gpioreg.ByName(shared.NtoPinName(pin)),
	}
}

func (g *GPIO) Init() error {
	if g == gpio.INVALID {
		return fmt.Errorf("pin %s is invalid", g.Name())
	}

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
