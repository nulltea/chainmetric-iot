package modules

import (
	"context"
	"sync"
	"time"

	"github.com/timoth-y/chainmetric-iot/controllers/device"
	"github.com/timoth-y/chainmetric-iot/model/events"
	"github.com/timoth-y/chainmetric-iot/shared"
	"github.com/timoth-y/go-eventdriver"
)

// moduleBase implements base functionality of the device.Module.
type moduleBase struct {
	*device.Device
	*sync.Once

	mid string
}

// withModuleBase can be used to embed moduleBase to device.Module implementation.
func withModuleBase(mid string) moduleBase {
	return moduleBase{
		Once: &sync.Once{},
		mid: mid,
	}
}

func (m *moduleBase) MID() string {
	return m.mid
}

func (m *moduleBase) Setup(device *device.Device) error {
	m.Device = device

	return nil
}

func (m *moduleBase) IsReady() bool {
	return m.Device != nil
}

func (m *moduleBase) Close() error {
	m.Device = nil

	return nil
}

// trySyncWithDeviceLifecycle attempts to synchronize device.Module execution with LifecycleManager module
// by waiting for device to login on network for a limited period of time,
// otherwise it will subscribe to events.DeviceLoggedOnNetwork and perform `postponedStart` when received one.
func (m *moduleBase) trySyncWithDeviceLifecycle(
	ctx context.Context,
	postponedStart func(context.Context),
) (synced bool) {
	if !m.waitUntilDeviceLogged() {
		m.Once = &sync.Once{}
		var cancel context.CancelFunc
		cancel = eventdriver.SubscribeHandler(events.DeviceLoggedOnNetwork, func(_ context.Context, _ interface{}) error {
			postponedStart(ctx)
			cancel()
			return nil
		})

		shared.Logger.Infof("device.Module '%s' is awaiting notification for the device login", m.MID())
		return false
	}

	return true
}

// waitUntilDeviceLogged checks whether the device.Device is logged on network with a specific intervals.
func (m *moduleBase) waitUntilDeviceLogged() bool {
	var attempts = 10
	for attempts > 0 {
		if m.IsLoggedToNetwork() {
			return true
		}

		time.Sleep(500 * time.Millisecond)
		attempts--
	}

	return false
}

// waitUntilSensorsDetected checks whether the sensors detected with a specific intervals.
func (m *moduleBase) waitUntilSensorsDetected() bool {
	var attempts = 5
	for attempts > 0 {
		if m.RegisteredSensors().NotEmpty() {
			return true
		}

		time.Sleep(250 * time.Millisecond)
		attempts--
	}

	return false
}
