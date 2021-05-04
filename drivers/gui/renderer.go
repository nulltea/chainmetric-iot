package gui

import (
	"fmt"
	"math"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/skip2/go-qrcode"
	"golang.org/x/image/font/gofont/gomedium"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/display"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

var (
	dev  display.Display
	width int
	height int
	ctx *gg.Context
)

func Init(display display.Display) {
	dev = display
	initContext()
}

func Text(text string) {
	clearFrame()

	font, err := truetype.Parse(gomedium.TTF)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(font, &truetype.Options{
		Size: 18,
	})
	ctx.SetFontFace(face)
	tw, th := ctx.MeasureString(text)
	ctx.DrawString(text, float64(height / 2) - tw / 2, float64(width / 2) + th / 2)
	ctx.Fill()

	dev.DrawAndRefresh(ctx.Image())
}

func TextWithIcon(text, icon string) {
	clearFrame()

	font, err := truetype.Parse(gomedium.TTF)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(font, &truetype.Options{
		Size: 18,
	})
	ctx.SetFontFace(face)
	tw, th := ctx.MeasureString(text)
	ctx.DrawString(text, float64(height / 2) - tw / 2, float64(width / 2) - th)
	iconImg, err := gg.LoadPNG(fmt.Sprintf("%s.png", icon))
	if err != nil {
		return
	}
	var (
		x = int(math.Round(float64(height / 2) - float64(iconImg.Bounds().Dy()) / 2))
		y = int(math.Round(float64(width / 2)))
	)

	ctx.DrawImage(iconImg, x, y)

	dev.DrawAndRefresh(ctx.Image())
}

func Success(msg string) {
	TextWithIcon(msg, "success")
}

func Warning(msg string) {
	TextWithIcon(msg, "warning")
}

func Error(msg string) {
	TextWithIcon(msg, "error")
}

func QR(data string) {
	qr, err := qrcode.New(data, qrcode.Medium); if err != nil {
		shared.Logger.Error("failed to create QR code image")
	}

	var (
		qrImg = qr.Image(width)
		x = int(math.Round(float64(height / 2) - float64(qrImg.Bounds().Dy()) / 2))
		y = int(math.Round(float64(width / 2) - float64(qrImg.Bounds().Dx()) / 2))
	)

	ctx.DrawImage(qrImg, x, y)
	dev.DrawAndRefresh(ctx.Image())
}

func Available() bool {
	return dev != nil && dev.Active()
}

func initContext() {
	width = dev.Bounds().Dx()
	height = dev.Bounds().Dy()
	ctx = gg.NewContext(width, height)

	ctx.Rotate(gg.Radians(90))
	ctx.Translate(0.0, -float64(height/2))
}

func clearFrame() {
	ctx.SetRGB(1, 1, 1)
	ctx.Clear()
	ctx.SetRGB(0, 0, 0)
}
