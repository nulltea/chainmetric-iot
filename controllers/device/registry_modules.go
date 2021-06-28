package device

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-iot/shared"
)

// ModulesRegistry defines pool of registered logical Module's extending the Device functionality.
type ModulesRegistry []Module

// Setup registers all logical device.Module's presented in ModulesRegistry onto the device.Device instance.
func (r ModulesRegistry) Setup(device *Device) {
	shared.Logger.Info("Setting up logical modules for the device...")

	for _, module := range r {
		if err := module.Setup(device); err == nil {
			shared.Logger.Info(
				"\033[32m[✔]\033[0m",
				fmt.Sprintf("Module '%s' setup compete", module.MID()),
			)
		} else {
			shared.Logger.Error(
				"\033[31m[✖]\033[0m",
				fmt.Sprintf("Failed to setup module '%s': %s", module.MID(), err),
			)
		}
	}
}

// Start starts all presented in ModulesRegistry logical device.Module's operational routine.
func (r ModulesRegistry) Start(ctx context.Context) {
	shared.Logger.Info("Device startup sequence started...")

	for _, m := range r {
		if m.IsReady() {
			m.Start(ctx)
			shared.Logger.Infof("\u001B[32m[⬤]\u001B[0m Module '%s' stated", m.MID())

			continue
		}

		shared.Logger.Warningf("\033[33m[🡆]\u001B[0m Module '%s' started is skipped due not readiness", m.MID())
	}

	shared.Logger.Info("Device is ready and running")
}

func (r ModulesRegistry) Close() {
	shared.Logger.Info("Device shutdown sequence started...")

	for _, m := range r {
		if !m.IsReady() {
			continue
		}

		if err := m.Close(); err != nil {
			shared.Logger.Error(errors.Wrapf(err, "failed to close '%s' module", m.MID()))
			continue
		}

		shared.Logger.Debugf("Module '%s' closed", m.MID())
	}

	shared.Logger.Info("Device has been shutdown")
}
