package peripheries

import (
	"context"
	"time"

	"github.com/go-ble/ble"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// Bluetooth defines BLE peripheral interface.
type Bluetooth struct {
	dev ble.Device
	srv *ble.Service

	name         string
	scanDuration time.Duration
	advDuration  time.Duration
}

// NewBluetooth creates new Bluetooth driver instance.
func NewBluetooth() *Bluetooth {
	return &Bluetooth{
		dev: shared.BluetoothDevice,
		srv: ble.NewService(ble.MustParse("1800")), // Generic access

		name: viper.GetString("bluetooth.name"),
		scanDuration: viper.GetDuration("bluetooth.scan_duration"),
		advDuration: viper.GetDuration("bluetooth.advertise_duration"),
	}
}

// Init performs Bluetooth driver initialisation.
func (b *Bluetooth) Init() (err error) {
	if err = b.dev.AddService(b.srv); err != nil {
		return errors.Wrap(err, "failed to enable bluetooth adapter")
	}

	return nil
}

// Scan performs scan for currently advertising bluetooth device.
func (b *Bluetooth) Scan() {
	var (
		ctx = ble.WithSigHandler(context.WithTimeout(context.Background(), b.scanDuration))
	)

	b.dev.Scan(ctx, false, func(adv ble.Advertisement) {

	})
}

// Advertise advertises device with previously configured name.
func (b *Bluetooth) Advertise() error {
	var (
		ctx = ble.WithSigHandler(context.WithTimeout(context.Background(), b.advDuration))
	)

	return b.dev.AdvertiseNameAndServices(ctx, b.name)
}

// Close closes Bluetooth connection and clears allocated resources.
func (b *Bluetooth) Close() error {
	return ble.Stop()
}
