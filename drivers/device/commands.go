package device

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-sensorsys/model/state"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

func (d *Device) Init() error {
	var (
		err error
		contract = d.client.Contracts.Devices
	)

	d.specs, err = d.DiscoverSpecs(true); if err != nil {
		return err
	}

	defer d.initHotswap()

	defer func() {
		go d.tryRepostCachedReadings()
	}()

	defer func() {
		d.active = err == nil && d.model != nil
	}()

	if id, is := isRegistered(); is {
		if d.model, _ = contract.Retrieve(id); d.model != nil {
			shared.Logger.Infof("Device specs has being updated in blockchain with id: %s", id)

			return contract.UpdateSpecs(id, d.specs)
		}

		shared.Logger.Warning("Device was removed from network, must re-initialize now")
	}

	if d.DisplayAvailable() {
		qr, err := qrcode.New(d.specs.Encode(), qrcode.Medium); if err != nil {
			return err
		}

		defer d.display.ClearAndRefresh()

		d.display.DrawAndRefresh(qr.Image(viper.GetInt("display.image_size")))
	} else {
		qrcode.WriteFile(d.specs.Encode(), qrcode.Medium, 280, "local/qr.png")
	}

	ctx, cancel := context.WithTimeout(d.ctx, viper.GetDuration("device.register_timeout_duration"))

	if err := contract.Subscribe(ctx, "inserted", func(device *models.Device, _ string) error {
		if device.Hostname == d.specs.Hostname {
			defer cancel()

			if err := storeIdentity(device.ID); err != nil {
				shared.Logger.Fatal(errors.Wrap(err, "failed to store device's identity file"))
			}

			shared.Logger.Infof("Device has being registered with id: %s", device.ID)
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

	if err := d.client.Contracts.Devices.Unbind(id); err != nil {
		return err
	}

	if err := os.Remove(viper.GetString("device.id_file_path")); err != nil {
		return errors.Wrap(err, "failed to remove device's identity file")
	}

	shared.Logger.Info("Device is been reset.")

	return nil
}

func (d *Device) NotifyOff() error {
	if d.client == nil || d.model == nil {
		return nil
	}
	return d.client.Contracts.Devices.UpdateState(d.model.ID, state.Offline)
}

func isRegistered() (string, bool) {
	id, err := ioutil.ReadFile(viper.GetString("device.id_file_path")); if err != nil {
		if os.IsNotExist(err) {
			return "", false
		}

		shared.Logger.Fatal(errors.Wrap(err, "failed to read device identity file"))
	}

	return string(id), true
}

func storeIdentity(id string) error {
	f, err := os.Create(viper.GetString("device.id_file_path")); if err != nil {
		return err
	}

	if _, err := f.WriteString(id); err != nil {
		return err
	}

	return nil
}
