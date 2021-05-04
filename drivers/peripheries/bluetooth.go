package peripheries

import (
	"github.com/pkg/errors"
	"tinygo.org/x/bluetooth"
)

// Bluetooth defines BLE peripheral interface.
type Bluetooth struct {
	*bluetooth.Adapter
	*bluetooth.Advertisement
}

// NewBluetooth creates new Bluetooth driver instance.
func NewBluetooth() *Bluetooth {
	return &Bluetooth{
		Adapter: bluetooth.DefaultAdapter,
	}
}

// Init performs Bluetooth driver initialisation.
func (b *Bluetooth) Init() error {
	if err := b.Enable(); err != nil {
		return errors.Wrap(err, "failed to enable bluetooth adapter")
	}

	b.Advertisement = b.DefaultAdvertisement()

	b.Advertisement.Configure(bluetooth.AdvertisementOptions{
		LocalName: "ChainMetric",
	})

	return nil
}
