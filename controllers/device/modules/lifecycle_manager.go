package modules

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-sensorsys/controllers/device"
	"github.com/timoth-y/chainmetric-sensorsys/controllers/gui"
	"github.com/timoth-y/chainmetric-sensorsys/core/net"
	"github.com/timoth-y/chainmetric-sensorsys/model"
	"github.com/timoth-y/chainmetric-sensorsys/model/events"
	"github.com/timoth-y/chainmetric-sensorsys/network/blockchain"
	"github.com/timoth-y/chainmetric-sensorsys/network/localnet"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
	"github.com/timoth-y/go-eventdriver"
)

// LifecycleManager defines device.Module for device.Device lifecycle managing.
type LifecycleManager struct {
	moduleBase
}

// WithLifecycleManager can be used to setup LifecycleManager logical device.Module onto the device.Device.
func WithLifecycleManager() device.Module {
	return &LifecycleManager{
		moduleBase: withModuleBase("LIFECYCLE_MANAGER"),
	}
}

func (m *LifecycleManager) Setup(device *device.Device) error {
	var (
		deviceName = viper.GetString("bluetooth.device_name")
	)

	if err := m.moduleBase.Setup(device); err != nil {
		return err
	}

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

func (m *LifecycleManager) Close() error {
	if err := m.notifyOff(); err != nil {
		return errors.Wrap(err, "failed to notify network about device shutdown")
	}

	return m.moduleBase.Close()
}

func (m *LifecycleManager) logInNetwork(ctx context.Context, id string) {
	if d, _ := blockchain.Contracts.Devices.Retrieve(id); d != nil {
		m.UpdateDeviceModel(d)
		eventdriver.EmitEvent(ctx, events.DeviceLoggedOnNetwork, nil)

		specs, err := m.discoverDeviceSpecs()
		if err != nil {
			shared.Logger.Error(errors.Wrap(err, "failed to discover device specs"))
			shared.Logger.Warning("Device specs update on network skipped")

			return
		}

		specs.State = models.DeviceOnline

		defer shared.MustExecute(func() error {
			return m.SetSpecs(func(ds *model.DeviceSpecs) {
				*ds = *specs
			})
		}, "failed to update device specs")

		shared.Logger.Infof("Device specs has being updated in blockchain with id: %s", id)
		return
	}

	shared.Logger.Warning("Device was removed from network, must re-initialize now")
	m.proceedToDeviceRegistration(ctx)
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

	specs.State = models.DeviceOnline

	ctx, cancel := context.WithTimeout(ctx, viper.GetDuration("device.register_timeout_duration"))

	// Try to start bluetooth advertisement:
	go func(ctx context.Context) {
		if err := localnet.Pair(ctx); err != nil {
			shared.Logger.Warning(errors.Wrap(err, "failed to advertise device via bluetooth"))
		}
	}(ctx)

	// Display registration payload as QR code:
	if gui.Available() {
		shared.Logger.Debug("Rendering QR")
		gui.RenderQRCode(specs.Encode())
	} else {
		// Alternative way to display registration QR code on Windows for debug purposes:
		_ = qrcode.WriteFile(specs.Encode(), qrcode.Medium, 320, "qr.png")
	}

	if err := contract.Subscribe(ctx, "inserted", func(dev *models.Device, _ string) error {
		if dev.Hostname == specs.Hostname {
			defer cancel()

			if err := m.storeIdentity(dev.ID); err != nil {
				shared.Logger.Fatal(errors.Wrap(err, "failed to store device's identity file"))
			}

			shared.Logger.Infof("Device has being registered with id: %s", dev.ID)
			m.UpdateDeviceModel(dev)
			eventdriver.EmitEvent(ctx, events.DeviceLoggedOnNetwork, nil)

			defer shared.MustExecute(func() error {
				return m.SetSpecs(func(ds *model.DeviceSpecs) {
					*ds = *specs
				})
			}, "failed to update device specs")

			defer gui.RenderSuccessMsg("Registration completed!")

			return nil
		}
		return nil
	}); err != nil {
		shared.Logger.Error(errors.Wrap(err, "failed to register device"))
		defer gui.RenderErrorMsg("Registration failed!")
	}
}

func (m *LifecycleManager) discoverDeviceSpecs() (*model.DeviceSpecs, error) {
	netEnv, err := net.GetNetworkEnvironmentInfo(); if err != nil {
		return nil, errors.Wrap(err, "failed to get network info")
	}

	// Attempting to get available sensor for engine to work with:
	var attempts = 10
	for attempts > 0 {
		if m.RegisteredSensors().NotEmpty() {
			break
		}

		time.Sleep(250 * time.Millisecond)
		attempts--
	}

	return &model.DeviceSpecs{
		Network: *netEnv,
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

func (m *LifecycleManager) notifyOff() error {
	if !m.IsLoggedToNetwork() {
		return nil
	}

	if err := m.SetState(models.DeviceOffline); err != nil {
		return err
	}

	return nil
}
