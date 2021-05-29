package modules

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-core/models"
	dev "github.com/timoth-y/chainmetric-sensorsys/drivers/device"
	"github.com/timoth-y/chainmetric-sensorsys/network/localnet"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// LocationManager implements Module for device.Device location management.
type LocationManager struct {
	*dev.Device
	*sync.Once
}


// WithLocationManager can be used to setup LocationManager logical Module onto the device.Device.
func WithLocationManager() Module {
	return &LocationManager{
		Once: &sync.Once{},
	}
}

func (m *LocationManager) MID() string {
	return "location_manager"
}

func (m *LocationManager) Setup(device *dev.Device) error {
	m.Device = device

	return nil
}

func (m *LocationManager) Start(ctx context.Context) {
	m.Do(func() {
		if err := localnet.Channels.Geo.Subscribe(ctx, func(location models.Location) error {
			if err := m.SetLocation(location); err != nil {
				return err
			}

			shared.Logger.Debugf("Device location was updated via Bluetooth tethering: %s", location.Name)

			return nil
		}); err != nil {
			shared.Logger.Error(errors.Wrap(err, "failed to subscribe to geo channel"))
		}
	})
}
