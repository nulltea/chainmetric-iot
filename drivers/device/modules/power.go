package modules

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"
	dev "github.com/timoth-y/chainmetric-sensorsys/drivers/device"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/power"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// PowerManager defines device.Device module for battery management.
type PowerManager struct {
	dev  *dev.Device
	ups  *power.UPSController
	once *sync.Once
}

// WithPowerManager can be used to setup PowerManager for the device.Device.
func WithPowerManager() Module {
	return &PowerManager{
		ups: power.NewUPSController(),
		once: &sync.Once{},
	}
}

func (m *PowerManager) Setup(device *dev.Device) error {
	if err := m.ups.Init(); err != nil {
		return errors.Wrap(err, "failed to initialize ups controller driver")
	}

	m.dev = device

	return nil
}

func (m *PowerManager) Start(ctx context.Context) {
	m.once.Do(func() {
		var (
			startTime  time.Time
			interval = viper.GetDuration("device.battery_check_interval")
		)

		go func() {
		LOOP:
			for {
				startTime = time.Now()

				level, err := m.ups.BatteryLevel()
				if err != nil {
					shared.Logger.Error(err)
				}

				plugged := m.ups.IsPlugged()

				if err = m.dev.SetBattery(models.DeviceBattery{
					Level: &level,
					PluggedIn: plugged,
				}); err != nil {
					shared.Logger.Error(err)
				}

				shared.Logger.Debugf("Device battery was updated: %d% (plugged: %s)", level, plugged)

				select {
				case <-time.After(interval - time.Since(startTime)):
				case <- ctx.Done():
					shared.Logger.Debug("Power management module routine ended.")
					break LOOP
				}
			}
		}()
	})
}
