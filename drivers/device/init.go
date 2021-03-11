package device

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/gateway/blockchain"
	"github.com/timoth-y/iot-blockchain-sensorsys/model"
	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

const (
	deviceIdentityFile = "device.id"
)

var (
	Specs *model.DeviceSpecs
)

func Init(dc *blockchain.DevicesContract) error {
	var err error

	Specs, err = DiscoverSpecs(); if err != nil {
		return err
	}

	if id, is := isRegistered(); is {
		if err := dc.UpdateSpecs(id, Specs); err != nil {
			return err
		}

		shared.Logger.Infof("Device specs been updated in blockchain with id: %s", id)

		return nil
	}

	// TODO: display QR code
	qrcode.Encode(Specs.Encode(), qrcode.Medium, 135)
	qrcode.WriteFile(Specs.EncodeJson(), qrcode.Medium, 135, "../qr.png")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Minute)

	if err := dc.Subscribe(ctx, "inserted", func(device models.Device) error {
		if device.Hostname == Specs.Hostname {
			defer cancel()

			if err := storeIdentity(device.ID); err != nil {
				shared.Logger.Fatal(errors.Wrap(err, "failed to store device's identity file"))
			}

			shared.Logger.Infof("Device is been registered with id: %s", device.ID)

			return dc.UpdateSpecs(device.ID, Specs)
		}

		return nil
	}); err != nil {
		return err
	}

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
