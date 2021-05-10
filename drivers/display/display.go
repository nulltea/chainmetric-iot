package display

import (
	"image"
)

// Display defines base methods for controlling display device.
type Display interface {
	// Init performs Display device initialization.
	Init() error
	// Sleep puts Display to deep-sleep mode to save power.
	Sleep() error
	// Reset performs hardware Reset of the Display.
	Reset() error
	// Draw sends image to Display. Use Refresh() or DrawAndRefresh() to display image.
	DrawImage(image.Image) error
	// Clear clears the Display's content.
	Clear() error
	// Bounds returns Display dimensions.
	Bounds() image.Rectangle
	// Active checks whether the Display device is connected and active.
	Active() bool
	// Close closes connection to Display device and clears allocated resources.
	Close() error
}
