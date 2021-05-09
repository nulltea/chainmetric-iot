package local

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripheries"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// Client defines the interface for low range network communication.
type (
	Client struct {
		dev *peripheries.Bluetooth

		Channels channels
	}

	channels struct {
		Geo *GeoLocationChannel
	}
)

// NewClient create new low range network communication Client.
func NewClient() *Client {
	return &Client{
		dev: peripheries.NewBluetooth(),

		Channels: channels{
			Geo: NewGeoLocationChannel(),
		},
	}
}

// Init performs initialisation of the Client.
func (c *Client) Init(name string) error {
	if !viper.GetBool("bluetooth.enabled") {
		return nil
	}

	if err := c.dev.Init(peripheries.WithDeviceName(name)); err != nil {
		return errors.Wrap(err, "failed to prepare Bluetooth driver")
	}

	if err := c.Channels.Geo.expose(c.dev); err != nil {
		return errors.Wrap(err, "failed to expose client to geo channel")
	}

	return nil
}

// Pair performs pairing via Bluetooth.
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
