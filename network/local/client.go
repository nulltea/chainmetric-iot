package local

import (
	"github.com/pkg/errors"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripheries"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

type Client struct {
	dev *peripheries.Bluetooth
}

func NewBluetoothClient() *Client {
	return &Client{
		dev: peripheries.NewBluetooth(),
	}
}

func (c *Client) Init() error {
	if err := c.dev.AddService(NewLocationService().Service); err != nil {
		return errors.Wrap(err, "failed to add location service to Bluetooth device")
	}

	return nil
}

func (c *Client) Pair() error {
	shared.Logger.Debug("Bluetooth pairing started")

	if err := c.dev.Advertise(); err != nil {
		return errors.Wrap(err, "failed to advertise device via Bluetooth")
	}

	return nil
}
