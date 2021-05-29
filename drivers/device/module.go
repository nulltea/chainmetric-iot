package device

import (
	"context"
)

// Module defines interface for Device extendable logical modules.
type Module interface {
	// MID returns Module ID.
	MID() string
	// Setup registers logical Module onto the Device instance.
	Setup(device *Device) error
	// IsReady determines whether the logical Module's Setup is complete and it is ready to Start.
	IsReady() bool
	// Start starts Module operational routine.
	Start(ctx context.Context)
	// Close stops Module gracefully and clears allocated resources.
	Close() error
}
