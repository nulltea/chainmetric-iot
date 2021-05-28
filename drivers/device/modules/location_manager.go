package modules

import (
	"context"
	"sync"

	"github.com/timoth-y/chainmetric-core/models"
	dev "github.com/timoth-y/chainmetric-sensorsys/drivers/device"
	"github.com/timoth-y/chainmetric-sensorsys/network/localnet"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// LocationManager defines device.Device module for location management.
type LocationManager struct {
	*dev.Device
	once *sync.Once
}

// WithLocationManager can be used to setup LocationManager module for the device.Device.
func WithLocationManager() Module {
	return &LocationManager{
		once: &sync.Once{},
	}
}

func (m *LocationManager) Setup(device *dev.Device) error {
	m.Device = device

	return nil
}

func (m *LocationManager) Start(ctx context.Context) {
	m.once.Do(func() {
		localnet.Channels.Geo.Subscribe(ctx, func(location models.Location) error {
			if err := m.SetLocation(location); err != nil {
				return err
			}

			shared.Logger.Debugf("Device location was updated via Bluetooth tethering: %s", location.Name)

			return nil
		})
	})
}
