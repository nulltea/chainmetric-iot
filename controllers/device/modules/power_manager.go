package modules

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-sensorsys/controllers/device"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/power"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// PowerManager implements device.Module for device.Device battery management.
type PowerManager struct {
	moduleBase

	ups  *power.UPSController
}

// WithPowerManager can be used to setup PowerManager logical device.Module onto the device.Device.
func WithPowerManager() device.Module {
	return &PowerManager{
		moduleBase: withModuleBase("POWER_MANAGER"),
		ups: power.NewUPSController(),
	}
}


func (m *PowerManager) Setup(device *device.Device) error {
	if err := m.ups.Init(); err != nil {
		return errors.Wrap(err, "failed to initialize ups controller driver")
	}

	return m.moduleBase.Setup(device)
}

func (m *PowerManager) Start(ctx context.Context) {
	go m.Do(func() {
		var (
			startTime  time.Time
			interval = viper.GetDuration("device.battery_check_interval")
		)

	LOOP:
		for {
			select {
			case <-time.After(interval - time.Since(startTime)):
			case <- ctx.Done():
				shared.Logger.Debug("Power management module routine ended.")
				break LOOP
			}

			startTime = time.Now()

			level, err := m.ups.BatteryLevel()
			if err != nil {
				shared.Logger.Error(err)
				continue
			}

			if level == 0 {
				continue
			} // Fuel gauge chip is not yet ready.

			plugged := m.ups.IsPlugged()

			if err = m.SetBattery(models.DeviceBattery{
				Level: &level,
				PluggedIn: plugged,
			}); err != nil {
				shared.Logger.Error(err)
			}

			shared.Logger.Debugf("Device battery: %d%% left (plugged: %t)", level, plugged)
		}
	})
}
