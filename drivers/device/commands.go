package device

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/models/requests"
	"github.com/timoth-y/chainmetric-core/utils"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/gui"
	"github.com/timoth-y/chainmetric-sensorsys/model/state"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

func (d *Device) Init() error {
	var (
		err error
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

	if err = d.handleDeviceRegistration(); err != nil {
		return err
	}

	if err = d.localnet.Init(d.model.Name); err != nil {
		return errors.Wrap(err, "failed to initialise local network client")
	}

	return nil
}

func (d *Device) ListenRemoteCommands() error {
	var (
		contract = d.client.Contracts.Devices
	)

	if !d.active {
		return nil
	}

	if d.model == nil || len(d.model.ID) == 0 {
		return errors.New("cannot listen to commands before device registration")
	}

	go func() {
		if err := contract.ListenCommands(d.ctx, d.model.ID, func(id string, cmd models.DeviceCommand, args ...interface{}) error {
			switch cmd {
			case models.DevicePauseCmd:
			case models.DeviceResumeCmd:
			case models.DevicePairingCmd:
				d.handleBluetoothPairingCmd(id)
			default:
				shared.Logger.Error(errors.Errorf("command '%s' is not supported", cmd))
			}
			return nil
		}); err != nil {
			shared.Logger.Error(err)
		}
	}()

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

func (d *Device) handleDeviceRegistration() error {
	var (
		contract = d.client.Contracts.Devices
	)

	if id, is := isRegistered(); is {
		if d.model, _ = contract.Retrieve(id); d.model != nil {
			shared.Logger.Infof("Device specs has being updated in blockchain with id: %s", id)

			return contract.UpdateSpecs(id, d.specs)
		}

		shared.Logger.Warning("Device was removed from network, must re-initialize now")
	}

	if gui.Available() {
		gui.RenderQRCode(d.specs.Encode())
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


func (d *Device) handleBluetoothPairingCmd(cmdID string) {
	var (
		contract = d.client.Contracts.Devices
		results = requests.DeviceCommandResultsSubmitRequest{
			Status: models.DeviceCmdCompleted,
		}
	)

	if err := d.localnet.Pair(); err != nil {
		results.Status = models.DeviceCmdFailed
		results.Error = utils.StringPointer(err.Error())
		shared.Logger.Error(err)
	}

	results.Timestamp = time.Now().UTC()

	if err := contract.SubmitCommandResults(cmdID, results); err != nil {
		shared.Logger.Error(err)
	}
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
