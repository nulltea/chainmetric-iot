package modules

import (
	"context"

	"github.com/pkg/errors"
	dev "github.com/timoth-y/chainmetric-sensorsys/drivers/device"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// Module defines interface for device.Device extendable logical modules.
type Module interface {
	// MID returns Module ID.
	MID() string
	// Setup registers logical Module onto the device.Device instance.
	Setup(device *dev.Device) error
	// Start starts Module operational routine.
	Start(ctx context.Context)
}

// Registry defines pool of registered logical Module's extending the Device functionality.
type Registry []Module

// Setup registers all logical Module's presented in Registry onto the device.Device instance.
func (r Registry) Setup(device *dev.Device) {
	// TODO: log modules setting up

	for _, module := range r {
		if err := module.Setup(device); err != nil {
			shared.Logger.Error(errors.Wrapf(err, "failed to setup '%s' module to device", module.MID()))
		}
	}
}

// Start starts all presented in Registry logical Module's operational routine.
func (r Registry) Start(ctx context.Context) {
	for _, module := range r {
		module.Start(ctx)
	}
}
