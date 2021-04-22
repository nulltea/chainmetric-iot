package peripherals

import (
	"fmt"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"

	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

// GPIO provides wrapper for GPIO peripheral
type GPIO struct {
	pin int
	gpio.PinIO
}

func NewGPIO(pin int) *GPIO {
	return &GPIO{
		pin: pin,
	}
}

func (g *GPIO) Init() error {
	pin := gpioreg.ByName(shared.NtoPinName(g.pin))

	if pin == gpio.INVALID || pin == nil {
		return fmt.Errorf("pin %s is invalid", pin.Name())
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
