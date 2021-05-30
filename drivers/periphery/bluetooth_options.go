package periphery

import (
	"time"

	"github.com/go-ble/ble"
)

// An BluetoothOption configures a Bluetooth driver.
type BluetoothOption interface {
	Apply(d *Bluetooth)
}

// BluetoothOptionFunc is a function that configures a Bluetooth driver.
type BluetoothOptionFunc func(d *Bluetooth)

// Apply calls BluetoothOptionFunc on the driver instance.
func (f BluetoothOptionFunc) Apply(d *Bluetooth) {
	f(d)
}

// WithDeviceName can be used to specify Bluetooth device identifier name.
// Default is the one specified in the configuration: viper.GetString("bluetooth.device_name").
func WithDeviceName(name string) BluetoothOption {
	return BluetoothOptionFunc(func(d *Bluetooth) {
		d.name = name
	})
}

// WithScanDuration can be used to specify timeout for Bluetooth scanning.
// Default is the one specified in the configuration: viper.GetString("bluetooth.scan_duration").
func WithScanDuration(du time.Duration) BluetoothOption {
	return BluetoothOptionFunc(func(d *Bluetooth) {
		d.scanDuration = du
	})
}

// WithAdvertisementDuration can be used to specify timeout for Bluetooth advertisement.
// Default is the one specified in the configuration: viper.GetString("bluetooth.advertise_duration").
func WithAdvertisementDuration(du time.Duration) BluetoothOption {
	return BluetoothOptionFunc(func(d *Bluetooth) {
		d.advDuration = du
	})
}

// WithAdvertisementServices can be used to specify service to expose during Bluetooth advertisement.
// Default is no service being exposed.
func WithAdvertisementServices(uuids ...ble.UUID) BluetoothOption {
	return BluetoothOptionFunc(func(d *Bluetooth) {
		d.advServices = append(d.advServices, uuids...)
	})
}
