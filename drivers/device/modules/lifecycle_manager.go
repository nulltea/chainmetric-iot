package modules

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"
	dev "github.com/timoth-y/chainmetric-sensorsys/drivers/device"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/gui"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/network"
	"github.com/timoth-y/chainmetric-sensorsys/model"
	"github.com/timoth-y/chainmetric-sensorsys/model/events"
	"github.com/timoth-y/chainmetric-sensorsys/network/blockchain"
	"github.com/timoth-y/chainmetric-sensorsys/network/localnet"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
	"github.com/timoth-y/go-eventdriver"
)

// LifecycleManager defines Module for device.Device lifecycle managing.
type LifecycleManager struct {
	*dev.Device
	*sync.Once
}

// WithLifecycleManager can be used to setup LifecycleManager logical Module onto the device.Device.
func WithLifecycleManager() Module {
	return &LifecycleManager{
		Once: &sync.Once{},
	}
}

func (m *LifecycleManager) MID() string {
	return "lifecycle-manager"
}

func (m *LifecycleManager) Setup(device *dev.Device) error {
	m.Device = device

	var (
		deviceName = viper.GetString("bluetooth.device_name")
	)

	if len(m.Name()) != 0 {
		deviceName = fmt.Sprintf("%s.%s", deviceName, m.Name())
	}

	if err := localnet.Init(deviceName); err != nil {
		shared.Logger.Warning(errors.Wrap(err, "failed to init localnet client"))
	}

	return nil
}

func (m *LifecycleManager) Start(ctx context.Context) {
	go m.Do(func() {
		if id, is := isRegistered(); is {
			m.logInNetwork(ctx, id)
		} else {
			m.proceedToDeviceRegistration(ctx)
		}

		eventdriver.SubscribeHandler(events.DeviceRemovedFromNetwork, func(_ context.Context, _ interface{}) error {
			return errors.Wrap(m.resetDevice(true), "failed to reset device")
		})
	})
}

func (m *LifecycleManager) logInNetwork(ctx context.Context, id string) {
	if d, _ := blockchain.Contracts.Devices.Retrieve(m.ID()); d != nil {
		m.UpdateDeviceModel(d)
		eventdriver.EmitEvent(ctx, events.DeviceLoggedOnNetwork, nil)

		specs, err := m.discoverDeviceSpecs()
		if err != nil {
			shared.Logger.Error(errors.Wrap(err, "failed to discover device specs"))
			shared.Logger.Warning("Device specs update on network skipped")

			return
		}

		if err := m.SetSpecs(func(dc *model.DeviceSpecs) {
			dc = specs
		}); err != nil {
			shared.Logger.Error(err)
			return
		}

		shared.Logger.Infof("Device specs has being updated in blockchain with id: %s", id)
		return
	}

	shared.Logger.Warning("Device was removed from network, must re-initialize now")
}

func (m *LifecycleManager) proceedToDeviceRegistration(ctx context.Context) {
	var (
		contract = blockchain.Contracts.Devices
		specs, err = m.discoverDeviceSpecs()
	)

	 if err != nil {
	 	shared.Logger.Error(errors.Wrap(err, "failed to discover device specs"))
	 	shared.Logger.Warning("Device registration couldn't proceed further without specification discovered")

		return
	 }

	ctx, cancel := context.WithTimeout(ctx, viper.GetDuration("device.register_timeout_duration"))

	// Try to start bluetooth advertisement:
	if err := localnet.Pair(ctx); err != nil {
		shared.Logger.Warning(errors.Wrap(err, "failed to advertise device via bluetooth"))
	}

	// Display registration payload as QR code:
	if gui.Available() {
		gui.RenderQRCode(m.Specs().Encode())
	} else {
		// Alternative way to display registration QR code on Windows for debug purposes:
		_ = qrcode.WriteFile(m.Specs().Encode(), qrcode.Medium, 320, "qr.png")
	}

	if err := contract.Subscribe(ctx, "inserted", func(dev *models.Device, _ string) error {
		if dev.Hostname == m.Specs().Hostname {
			defer cancel()

			if err := m.storeIdentity(dev.ID); err != nil {
				shared.Logger.Fatal(errors.Wrap(err, "failed to store device's identity file"))
			}

			shared.Logger.Infof("Device has being registered with id: %s", dev.ID)
			m.UpdateDeviceModel(dev)
			eventdriver.EmitEvent(ctx, events.DeviceLoggedOnNetwork, nil)

			return m.SetSpecs(func(ds *model.DeviceSpecs) {
				ds = specs
			})
		}
		return nil
	}); err != nil {
		shared.Logger.Error(errors.Wrap(err, "failed to register device"))
	}
}

func (m *LifecycleManager) discoverDeviceSpecs() (*model.DeviceSpecs, error) {
	net, err := network.GetNetworkEnvironmentInfo(); if err != nil {
		return nil, errors.Wrap(err, "failed to get network info")
	}

	// Attempting to get available sensor for engine to work with:
	var attempts = 5
	for attempts > 0 {
		if m.RegisteredSensors().NotEmpty() {
			break
		}

		time.Sleep(250 * time.Millisecond)
		attempts--
	}

	return &model.DeviceSpecs{
		Network: *net,
		Supports: m.RegisteredSensors().SupportedMetrics(),
	}, nil
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

func (m *LifecycleManager) storeIdentity(id string) error {
	f, err := os.Create(viper.GetString("device.id_file_path")); if err != nil {
		return err
	}

	if _, err := f.WriteString(id); err != nil {
		return err
	}

	return nil
}

// resetDevice resets device.Device by removing stored identity and all allocated resources.
// Use `forceful` to specify that device must be reset since it has been removed from network.
func (m *LifecycleManager) resetDevice(forceful bool) error {
	if !forceful {
		id, is := isRegistered(); if !is {
			return nil
		}

		if err := blockchain.Contracts.Devices.Unbind(id); err != nil {
			return err
		}
	}

	if err := os.Remove(viper.GetString("device.id_file_path")); err != nil {
		return errors.Wrap(err, "failed to remove device's identity file")
	}

	shared.Logger.Info("Device has been reset")

	return nil
}

func (m *LifecycleManager) notifyOff() {
	if !m.IsLoggedToNetwork() {
		return
	}

	if err := m.SetState(models.DeviceOffline); err != nil {
		shared.Logger.Error(err)
	}
}


// waitUntilDeviceLogged checks whether the device.Device is logged on network with a specific intervals.
func waitUntilDeviceLogged(d *dev.Device) bool {
	var attempts = 5
	for attempts > 0 {
		if d.RegisteredSensors().NotEmpty() {
			return true
		}

		time.Sleep(250 * time.Millisecond)
		attempts--
	}

	return false
}
