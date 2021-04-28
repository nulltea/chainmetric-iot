package display

import (
	"image"
)

// Display defines base methods for controlling display device.
type Display interface {
	Init() error
	Halt() error
	Sleep() error
	Reset()
	DrawImage(image.Image) error
	Bounds() image.Rectangle
	Active() bool
	Close() error
}
