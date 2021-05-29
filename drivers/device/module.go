package device

import (
	"context"
)

// Module defines interface for device.Device extendable logical modules.
type Module interface {
	// MID returns device.Module ID.
	MID() string
	// Setup registers logical device.Module onto the device.Device instance.
	Setup(device *Device) error
	// Start starts device.Module operational routine.
	Start(ctx context.Context)
}
