package modules

import (
	"context"

	dev "github.com/timoth-y/chainmetric-sensorsys/drivers/device"
)

// Module defines interface for device.Device extendable modules.
type Module interface {
	Setup(device *dev.Device) error
	Start(ctx context.Context)
}
