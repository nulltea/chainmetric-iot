package device

import (
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/display"
	"github.com/timoth-y/iot-blockchain-sensorsys/gateway/blockchain"
)

type Device struct {
	display display.Display
	client *blockchain.Client
}

func NewDevice() *Device {
	return &Device{

	}
}

func (d *Device) SetDisplay(dp display.Display) {
	d.display = dp
}

func (d *Device) SetClient(client *blockchain.Client) {
	d.client = client
}

func (d *Device) Close() error {
	if err := d.display.Close(); err != nil {
		return err
	}

	d.client.Close()

	return nil
}
