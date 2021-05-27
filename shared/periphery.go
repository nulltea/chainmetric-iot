package shared

import (
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"periph.io/x/periph/host"
)

var(
	// BluetoothDevice is a default Bluetooth device instance.
	BluetoothDevice ble.Device
)

// initPeriphery initialises peripheral drivers host (I2C, SPI, Bluetooth)
func initPeriphery() {
	var err error

	if !viper.GetBool("bluetooth.enabled") {
		return
	}

	if _, err = host.Init(); err != nil {
		Logger.Fatal(errors.Wrap(err, "failed to initialise peripheral host"))
	}

	if BluetoothDevice, err = linux.NewDeviceWithName(viper.GetString("bluetooth.name")); err != nil {
		Logger.Fatal(errors.Wrap(err, "failed to create bluetooth device"))
	}

	ble.SetDefaultDevice(BluetoothDevice)
}

func closePeriphery() {
	if BluetoothDevice != nil {
		Execute(BluetoothDevice.Stop, "failed to close bluetooth periphery")
	}
}
