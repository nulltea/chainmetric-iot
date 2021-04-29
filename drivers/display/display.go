package display

import (
	"image"
)

// Display defines base methods for controlling display device.
type Display interface {
	Init() error
	Sleep() error
	Reset() error
	Draw(image.Image) error
	DrawAndRefresh(image.Image) error
	Clear() error
	Refresh() error
	Bounds() image.Rectangle
	Active() bool
	Close() error
}
