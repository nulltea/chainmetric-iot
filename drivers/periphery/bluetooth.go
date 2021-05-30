package periphery

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
	advServices  []ble.UUID
}

// NewBluetooth creates new Bluetooth driver instance.
func NewBluetooth(options ...BluetoothOption) *Bluetooth {
	return (&Bluetooth{
		Device: shared.BluetoothDevice,
		name: viper.GetString("bluetooth.device_name"),
		scanDuration: viper.GetDuration("bluetooth.scan_duration"),
		advDuration: viper.GetDuration("bluetooth.advertise_duration"),
		advServices: []ble.UUID{},
	}).ApplyOptions(options...)
}

// Init performs initialisation of the Bluetooth driver.
func (b *Bluetooth) Init(options ...BluetoothOption) error {
	b.ApplyOptions(options...)

	return nil
}

// ApplyOptions applies Bluetooth configuration parameters by given `options`.
func (b *Bluetooth) ApplyOptions(options ...BluetoothOption) *Bluetooth {
	for i := range options {
		options[i].Apply(b)
	}

	return b
}

// Scan performs scan for currently advertising bluetooth device.
func (b *Bluetooth) Scan(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, b.scanDuration)
	ctx = ble.WithSigHandler(ctx, cancel)
	defer cancel()

	return b.Device.Scan(ctx, false, func(adv ble.Advertisement) {

	})
}

// Advertise advertises device with previously configured name.
func (b *Bluetooth) Advertise(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, b.scanDuration)
	ctx = ble.WithSigHandler(ctx, cancel)
	defer cancel()

	return b.Device.AdvertiseNameAndServices(ctx, b.name, b.advServices...)
}

// Close closes Bluetooth connection and clears allocated resources.
func (b *Bluetooth) Close() error {
	if b.Device == nil {
		return nil
	}

	return b.Device.Stop()
}

