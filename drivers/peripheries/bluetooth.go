package peripheries

import (
	"context"
	"time"

	"github.com/go-ble/ble"
	"github.com/spf13/viper"

	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// Bluetooth defines BLE peripheral interface.
type Bluetooth struct {
	ble.Device
	name         string
	scanDuration time.Duration
	advDuration  time.Duration
}

// NewBluetooth creates new Bluetooth driver instance.
func NewBluetooth() *Bluetooth {
	return &Bluetooth{
		Device: shared.BluetoothDevice,
		name: viper.GetString("bluetooth.device_name"),
		scanDuration: viper.GetDuration("bluetooth.scan_duration"),
		advDuration: viper.GetDuration("bluetooth.advertise_duration"),
	}
}


// Scan performs scan for currently advertising bluetooth device.
func (b *Bluetooth) Scan() {
	var (
		ctx = ble.WithSigHandler(context.WithTimeout(context.Background(), b.scanDuration))
	)

	b.Device.Scan(ctx, false, func(adv ble.Advertisement) {

	})
}

// Advertise advertises device with previously configured name.
func (b *Bluetooth) Advertise() error {
	var (
		ctx = ble.WithSigHandler(context.WithTimeout(context.Background(), b.advDuration))
	)

	return b.Device.AdvertiseNameAndServices(ctx, b.name)
}

// Close closes Bluetooth connection and clears allocated resources.
func (b *Bluetooth) Close() error {
	return b.Device.Stop()
}
