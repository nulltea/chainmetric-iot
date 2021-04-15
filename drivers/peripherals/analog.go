package peripherals

import (
	"errors"
	"fmt"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"

	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

type AnalogMCP3208 struct {
	spi *SPI
	pin string
	cs  gpio.PinOut
	tx  []byte
	rx  []byte
	CH0 ADCPin
	CH1 ADCPin
	CH2 ADCPin
	CH3 ADCPin
	CH4 ADCPin
	CH5 ADCPin
	CH6 ADCPin
	CH7 ADCPin
}

// ADCPin is the implementation of the ADConverter interface.
type ADCPin struct {
	gpio.PinOut
	dev *AnalogMCP3208
}

// NewAnalogMCP3208 returns a new MCP3008 driver.
func NewAnalogMCP3208(bus string, csPin int) *AnalogMCP3208 {
	return &AnalogMCP3208{
		spi: NewSPI(bus),
		pin: shared.NtoPinName(csPin),
		tx: make([]byte, 3),
		rx: make([]byte, 3),
	}
}

// Configure sets up the device for communication
func (d *AnalogMCP3208) Init() error {
	if d.cs = gpioreg.ByName(d.pin); d.cs == gpio.INVALID {
		return fmt.Errorf("failed to connect CS pin via %s", d.pin)
	}

	if err := d.spi.Init(); err != nil {
		return err
	}

	// setup all channels
	d.CH0 = d.GetADC(0)
	d.CH1 = d.GetADC(1)
	d.CH2 = d.GetADC(2)
	d.CH3 = d.GetADC(3)
	d.CH4 = d.GetADC(4)
	d.CH5 = d.GetADC(5)
	d.CH6 = d.GetADC(6)
	d.CH7 = d.GetADC(7)

	return nil
}

// Read analog data from channel
func (d *AnalogMCP3208) Read(ch int) (uint16, error) {
	if ch < 0 || ch > 7 {
		return 0, errors.New("invalid channel for MCP3208 Read")
	}

	return d.GetADC(ch).Get(), nil
}

// GetADC returns an ADC for a specific channel.
func (d *AnalogMCP3208) GetADC(ch int) ADCPin {
	return ADCPin{gpioreg.ByName(shared.NtoPinName(ch)), d}
}

// Get the current reading for a specific ADCPin.
func (p ADCPin) Get() uint16 {
	p.dev.tx[0] = 0x01
	p.dev.tx[1] = byte(8 + p.Number()) << 4
	p.dev.tx[2] = 0x00

	p.dev.cs.Out(gpio.Low)
	p.dev.spi.Tx(p.dev.tx, p.dev.rx)

	// scale result to 16bit value like other ADCs
	result := uint16(p.dev.rx[1] & 0x3) << 8 + uint16(p.dev.rx[2]) << 6
	p.dev.cs.Out(gpio.High)

	return result
}
