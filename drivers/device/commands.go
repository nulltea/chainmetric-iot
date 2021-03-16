package device

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

const (
	deviceIdentityFile = "device.id"
)

func (d *Device) Init() error {
	var (
		err error
		contract = d.client.Contracts.Devices
	)

	d.specs, err = d.DiscoverSpecs(); if err != nil {
		return err
	}

	if id, is := isRegistered(); is {
		if d.model, _ = contract.Retrieve(id); d.model == nil {
			if err := contract.UpdateSpecs(id, d.specs); err != nil {
				return err
			}

			shared.Logger.Infof("Device specs been updated in blockchain with id: %s", id)

			return nil
		}

		shared.Logger.Warning("device was removed from network, must re-initialize now")
	}

	if d.display != nil {
		qr, err := qrcode.New(d.specs.Encode(), qrcode.Medium); if err != nil {
			return err
		}

		d.display.PowerOn()
		defer d.display.PowerOff()

		d.display.DrawImage(qr.Image(d.config.Display.ImageSize))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Minute)

	if err := contract.Subscribe(ctx, "inserted", func(device *models.Device, _ string) error {
		if device.Hostname == d.specs.Hostname {
			defer cancel()

			if err := storeIdentity(device.ID); err != nil {
				shared.Logger.Fatal(errors.Wrap(err, "failed to store device's identity file"))
			}

			shared.Logger.Infof("Device is been registered with id: %s", device.ID)
			d.model = device
			return contract.UpdateSpecs(device.ID, d.specs)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (d *Device) Reset() error {
	id, is := isRegistered(); if !is {
		return nil
	}

	if err := d.client.Contracts.Devices.Remove(id); err != nil {
		return err
	}

	if err := os.Remove(deviceIdentityFile); err != nil {
		return errors.Wrap(err, "failed to remove device's identity file")
	}

	shared.Logger.Info("Device is been reset.")

	return nil
}


func isRegistered() (string, bool) {
	id, err := ioutil.ReadFile(deviceIdentityFile); if err != nil {
		if os.IsNotExist(err) {
			return "", false
		}

		shared.Logger.Fatal(errors.Wrap(err, "failed to read device identity file"))
	}

	return string(id), true
}

func storeIdentity(id string) error {
	f, err := os.Create(deviceIdentityFile); if err != nil {
		return err
	}

	if _, err := f.WriteString(id); err != nil {
		return err
	}

	return nil
}
