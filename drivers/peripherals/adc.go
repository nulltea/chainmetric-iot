package peripherals

import (
	"fmt"

	"github.com/pkg/errors"
)

type AnalogMCP3208 struct {
	spi *SPI
	cs  *GPIO
	channels []*AnalogChannel
}

// AnalogChannel is the implementation of the ADConverter interface.
type AnalogChannel struct {
	ch int
	dev *AnalogMCP3208
	tx  []byte
	rx  []byte
}

// NewAnalogMCP3208 returns a new MCP3008 driver.
func NewAnalogMCP3208(bus string, csPin int) *AnalogMCP3208 {
	return &AnalogMCP3208{
		spi: NewSPI(bus),
		cs: NewGPIO(csPin),
		channels: make([]*AnalogChannel, 8),
	}
}

// Init sets up the device for communication and prepares all available channels
func (d *AnalogMCP3208) Init() error {
	if err := d.cs.Init(); err != nil {
		return errors.Wrap(err, "failed to connect CS pin")
	}

	if err := d.spi.Init(); err != nil {
		return err
	}

	for ch := 0; ch < 8; ch++ {
		d.channels[ch] = d.GetChannel(ch)
	}

	return nil
}

// Read analog reading from a specified channel.
func (d *AnalogMCP3208) Read(ch int) (uint16, error) {
	if ch < 0 || ch > 7 {
		return 0, fmt.Errorf("channel CH%d is not sopported by MCP3208 device", ch)
	}

	return d.GetChannel(ch).Get(), nil
}

// GetChannel returns an AnalogChannel for a specified channel number.
func (d *AnalogMCP3208) GetChannel(ch int) *AnalogChannel {
	return &AnalogChannel{
		ch: ch,
		dev: d,
		tx: make([]byte, 3),
		rx: make([]byte, 3),
	}
}

// Get analog reading from current AnalogChannel.
func (p AnalogChannel) Get() uint16 {
	p.tx[0] = 0x01
	p.tx[1] = byte(8 + p.ch) << 4
	p.tx[2] = 0x00

	p.dev.cs.Low()
	p.dev.spi.Tx(p.tx, p.rx)

	// scale result to 16bit value
	result := uint16(p.rx[1] & 0x3) << 8 + uint16(p.rx[2]) << 6
	p.dev.cs.High()

	return result
}
