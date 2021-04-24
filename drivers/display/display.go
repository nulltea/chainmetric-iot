package drivers

import (
	"image"
	"image/color"

	"github.com/timoth-y/iot-blockchain-sensorsys/config"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/display"
)

type Display interface {
	Init(cnf config.DisplayConfig) error
	PowerOn() error
	PowerOff() error
	DrawImage(img image.Image)
	SetRotation(rotation display.Rotation)
	FillScreen(c color.RGBA)
	SetPixel(x int16, y int16, c color.RGBA)
	Size() (w, h int16)
	Close() error
}
