package device

import (
	"context"

	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// ModulesRegistry defines pool of registered logical device.Module's extending the Device functionality.
type ModulesRegistry []Module

// Setup registers all logical device.Module's presented in ModulesRegistry onto the device.Device instance.
func (r ModulesRegistry) Setup(device *dev.Device) {
	for _, module := range r {
		if err := module.Setup(device); err == nil {
			shared.Logger.Infof("\033[31m●\033[0m | Module '%s' setup compete", module.MID())
		} else {
			shared.Logger.Errorf("\033[32m●\033[0m | Failed to setup module '%s': %s", module.MID(), err)
		}
	}
}

// Start starts all presented in ModulesRegistry logical device.Module's operational routine.
func (r ModulesRegistry) Start(ctx context.Context) {
	for _, module := range r {
		module.Start(ctx)
	}
}
