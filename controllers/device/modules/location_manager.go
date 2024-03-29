package modules

import (
	"context"

	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-iot/controllers/device"
	"github.com/timoth-y/chainmetric-iot/model/events"
	"github.com/timoth-y/chainmetric-iot/network/localnet"
	"github.com/timoth-y/chainmetric-iot/shared"
	"github.com/timoth-y/go-eventdriver"
)

// LocationManager implements device.Module for device.Device location management.
type LocationManager struct {
	moduleBase
}


// WithLocationManager can be used to setup LocationManager logical device.Module onto the device.Device.
func WithLocationManager() device.Module {
	return &LocationManager{
		moduleBase: withModuleBase("LOCATION_MANAGER"),
	}
}

func (m *LocationManager) Start(ctx context.Context) {
	go m.Do(func() {
		if err := localnet.Channels.Geo.Subscribe(ctx, func(location models.Location) error {
			if err := m.SetLocation(location); err != nil {
				return err
			}

			eventdriver.EmitEvent(ctx, events.LocationUpdateReceived, location)

			shared.Logger.Debugf("Device location was updated via Bluetooth tethering: %s", location.Name)

			return nil
		}); err != nil {
			shared.Logger.Error(errors.Wrap(err, "failed to subscribe to geo channel"))
		}
	})
}
