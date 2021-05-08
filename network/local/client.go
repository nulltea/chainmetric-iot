package local

import (
	"github.com/pkg/errors"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripheries"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// Client defines the interface for low range network communication.
type Client struct {
	dev *peripheries.Bluetooth

	Topics topics
}

type topics struct {
	Location LocationTopic
}

// NewClient create new low range network communication client.
func NewClient() *Client {
	return &Client{
		dev: peripheries.NewBluetooth(),
	}
}

// Init performs initialisation of the Bluetooth Client.
func (c *Client) Init(name string) error {
	if err := c.dev.Init(peripheries.WithDeviceName(name)); err != nil {
		return errors.Wrap(err, "failed to prepare Bluetooth driver")
	}

	if err := c.dev.AddService(NewLocationTopic().Service); err != nil {
		return errors.Wrap(err, "failed to add location service to Bluetooth device")
	}

	return nil
}

// Pair performs bluetooth pairing.
func (c *Client) Pair() error {
	shared.Logger.Debug("Bluetooth pairing started")

	if err := c.dev.Advertise(); err != nil {
		return errors.Wrap(err, "failed to advertise device via Bluetooth")
	}

	return nil
}

// Close closes local network connection and clears allocated resources.
func (c *Client) Close() error {
	return c.dev.Close()
}
