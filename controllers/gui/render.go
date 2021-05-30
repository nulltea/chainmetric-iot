package gui

import (
	"fmt"
	"image"
	"math"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
	"github.com/timoth-y/chainmetric-sensorsys/core"
	"golang.org/x/image/font/gofont/gomedium"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

var (
	dev         core.Display
	frameWidth  int
	frameHeight int
	ctx         *gg.Context
)

// Init initialises GUI agent.
func Init(display core.Display) {
	dev = display
	initContext()
}

// RenderText  displays frame with `msg` text.
func RenderText(msg string) {
	initContext()

	font, err := truetype.Parse(gomedium.TTF)
	if err != nil {
		shared.Logger.Error(errors.Wrap(err, "failed to parse font"))
	}
	face := truetype.NewFace(font, &truetype.Options{
		Size: 18,
	})
	ctx.SetFontFace(face)
	tw, th := ctx.MeasureString(msg)
	ctx.DrawString(msg, float64(frameHeight/ 2) - tw / 2, float64(frameWidth/ 2) + th / 2)

	ShowFrame()
}

// RenderText displays frame with `msg` text and `icon` image.
func RenderTextWithIcon(text, icon string) {
	initContext()

	var (
		iconImg image.Image
		iconPath = iconPath(icon)
		font, err = truetype.Parse(goregular.TTF)
		face = truetype.NewFace(font, &truetype.Options{
			Size: 18,
		})
	)

	if err != nil {
		shared.Logger.Error(errors.Wrap(err, "failed to create font face"))
		return
	}

	var (
		tw, th = ctx.MeasureString(text)
		tx, ty = float64(frameHeight/ 2) - tw / 2, float64(frameWidth/ 2) - th
	)

	if iconImg, err = gg.LoadPNG(iconPath); err != nil {
		shared.Logger.Error(errors.Wrapf(err, "failed to load from path '%s'", iconPath))
		return
	}

	var (
		ib = iconImg.Bounds()
		ix = int(math.Round(float64(frameHeight/ 2) - float64(ib.Dy()) / 2))
		iy = int(math.Round(float64(frameWidth / 2)))
	)

	ctx.SetFontFace(face)
	ctx.DrawString(text, tx, ty)
	ctx.DrawImage(iconImg, ix, iy)

	ShowFrame()
}

// RenderSuccessMsg displays frame with `msg` text and "success" icon.
func RenderSuccessMsg(msg string) {
	RenderTextWithIcon(msg, "success")
}

// RenderWarningMsg displays frame with `msg` and "warning" icon.
func RenderWarningMsg(msg string) {
	RenderTextWithIcon(msg, "warning")
}

// RenderErrorMsg displays frame with `msg` text and "error" icon.
func RenderErrorMsg(msg string) {
	RenderTextWithIcon(msg, "error")
}

// RenderQRCode displays frame with QR code of `data`.
func RenderQRCode(data string) {
	qr, err := qrcode.New(data, qrcode.Medium); if err != nil {
		shared.Logger.Error(errors.Wrap(err, "failed to create QR code image"))
		return
	}

	var (
		qrImg = qr.Image(frameWidth)
		x = int(math.Round(float64(frameHeight/ 2) - float64(qrImg.Bounds().Dy()) / 2))
		y = int(math.Round(float64(frameWidth/ 2) - float64(qrImg.Bounds().Dx()) / 2))
	)

	ctx.DrawImage(qrImg, x, y)
	ctx.Fill()

	ShowFrame()
}

// ShowFrame displays frame with rendered context.
func ShowFrame() {
	ctx.Fill()
	dev.DrawAndRefresh(ctx.Image())
}

// Available checks whether the GUI is available.
func Available() bool {
	return dev != nil && dev.Active()
}

func initContext() {
	frameWidth = dev.Bounds().Dx()
	frameHeight = dev.Bounds().Dy()
	ctx = gg.NewContext(frameWidth, frameHeight)

	clearFrame()

	ctx.Rotate(gg.Radians(90))
	ctx.Translate(0.0, -float64(frameHeight/2))
}

func clearFrame() {
	ctx.SetRGB(1, 1, 1)
	ctx.Clear()
	ctx.SetRGB(0, 0, 0)
}

func iconPath(icon string) string {
	return fmt.Sprintf("drivers/gui/assets/%s.png", icon)
}
