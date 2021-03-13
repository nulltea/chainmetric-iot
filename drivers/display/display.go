package display

import (
	"image"
	"image/color"
)

type Display interface {
	PowerOn() error
	PowerOff() error
	DrawImage(img image.Image)
	SetRotation(rotation Rotation)
	FillScreen(c color.RGBA)
	SetPixel(x int16, y int16, c color.RGBA)
	Size() (w, h int16)
	Close() error
}
