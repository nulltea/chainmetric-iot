package modules

import (
	"context"
	"sync"

	"github.com/timoth-y/chainmetric-core/models"
	dev "github.com/timoth-y/chainmetric-sensorsys/drivers/device"
	"github.com/timoth-y/chainmetric-sensorsys/network/localnet"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// LocationManagement defines device.Device module for location management.
type LocationManagement struct {
	dev  *dev.Device
	once *sync.Once
}

// WithLocationManagement can be used to setup LocationManagement module for the device.Device.
func WithLocationManagement() Module {
	return &LocationManagement{
		once: &sync.Once{},
	}
}

func (m *LocationManagement) Setup(device *dev.Device) error {
	m.dev = device

	return nil
}

func (m *LocationManagement) Start(ctx context.Context) {
	m.once.Do(func() {
		localnet.Channels.Geo.Subscribe(ctx, func(location models.Location) error {
			if err := m.dev.SetLocation(location); err != nil {
				return err
			}

			shared.Logger.Debugf("Device location was updated via Bluetooth tethering: %s", location.Name)

			return nil
		})
	})
}
