package display

import (
	"image"
	"image/color"

	"github.com/timoth-y/chainmetric-sensorsys/model/config"
)

type Display interface {
	Init(cnf config.DisplayConfig) error
	PowerOn() error
	PowerOff() error
	DrawImage(img image.Image)
	SetRotation(rotation Rotation)
	FillScreen(c color.RGBA)
	SetPixel(x int16, y int16, c color.RGBA)
	Size() (w, h int16)
	Active() bool
	Close() error
}
