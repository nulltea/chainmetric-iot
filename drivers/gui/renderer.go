package gui

import (
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/gomedium"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/display"
)

var (
	dev  display.Display
)

func Init(display display.Display) {
	dev = display
}

func Text(text string) {
	bounds := dev.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	dc := gg.NewContext(w, h)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	dc.Rotate(gg.Radians(90))
	dc.Translate(0.0, -float64(h/2))
	font, err := truetype.Parse(gomedium.TTF)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(font, &truetype.Options{
		Size: 16,
	})
	dc.SetFontFace(face)
	tw, th := dc.MeasureString(text)
	padding := 8.0
	dc.DrawRoundedRectangle(padding*2, padding*2, tw+padding*2, th+padding, 10)
	dc.Stroke()
	dc.DrawString(text, padding*3, padding*2+th)
	for i := 0; i < 10; i++ {
		dc.DrawCircle(float64(30+(10*i)), 100, 5)
	}
	for i := 0; i < 10; i++ {
		dc.DrawRectangle(float64(30+(10*i)), 80, 5, 5)
	}
	dc.Fill()
	dev.DrawAndRefresh(dc.Image())
}
