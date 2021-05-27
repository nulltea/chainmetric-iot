package localnet

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripheries"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// Client defines the interface for low range network communication.
type Client struct {
	dev *peripheries.Bluetooth
}

var (
	client = &Client{
		dev: peripheries.NewBluetooth(),
	}

	// Channels exposes available channel for local network communication.
	Channels = struct {
		Geo *GeoLocationChannel
	} {
		Geo: NewGeoLocationChannel(),
	}
)

// Init performs initialisation of the Client.
func Init(name string) error {
	if !viper.GetBool("bluetooth.enabled") {
		return nil
	}

	if err := client.dev.Init(
		peripheries.WithDeviceName(name),
		peripheries.WithAdvertisementServices(Channels.Geo.uuid),
	); err != nil {
		return errors.Wrap(err, "failed to prepare Bluetooth driver")
	}

	if err := Channels.Geo.expose(client.dev); err != nil {
		return errors.Wrap(err, "failed to expose client to geo channel")
	}

	return nil
}

// Pair performs pairing via Bluetooth.
func Pair() error {
	shared.Logger.Debug("Bluetooth pairing started")

	if err := client.dev.Advertise(); err != nil {
		return errors.Wrap(err, "failed to advertise device via Bluetooth")
	}

	return nil
}

// Close closes local network connection and clears allocated resources.
func Close() error {
	return client.dev.Close()
}
