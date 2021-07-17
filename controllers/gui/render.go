package gui

import (
	"bytes"
	"fmt"
	"image"
	"math"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"
	core "github.com/timoth-y/chainmetric-iot/core/dev"
	"github.com/timoth-y/chainmetric-iot/shared"
	"github.com/wcharczuk/go-chart"
	"golang.org/x/image/font/gofont/gomedium"
	"golang.org/x/image/font/gofont/goregular"
)

var (
	dev         core.Display
	frameHeight int
	frameWidth  int
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

	var (
		vIndent = 0.0
	)

	for _, line := range strings.Split(msg, "\n") {
		tw, th := ctx.MeasureString(line)
		ctx.DrawString(line, float64(frameWidth/ 2) - tw / 2, float64(frameHeight/ 2) + th / 2 + vIndent)
		vIndent += th + 5
	}

	ShowFrame()
}

// RenderWithChart displays frame with chart and `msg` text.
func RenderWithChart(msg string, v ...float64) {
	initContext()

	font, err := truetype.Parse(gomedium.TTF)
	if err != nil {
		shared.Logger.Error(errors.Wrap(err, "failed to parse font"))
	}

	face := truetype.NewFace(font, &truetype.Options{
		Size: 14,
	})
	ctx.SetFontFace(face)

	var vIndent = 0.0

	for _, line := range strings.Split(msg, "\n") {
		_, th := ctx.MeasureString(line)
		ctx.DrawString(line, 0, 0 + th + vIndent)
		vIndent += th + 2
	}

	xValues := make([]float64, len(v))
	for i := range v {
		xValues[i] = float64(i) + 1
	}

	chart.DefaultFillColor = chart.ColorBlack
	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: xValues,
				YValues: v,
				Style: chart.Style{
					Show: true,
					StrokeWidth: 3,
					DotColor: chart.ColorBlack,
					StrokeColor: chart.ColorBlack,
				},
			},
		},
		Height: frameHeight / 2,
		Width: frameWidth,
	}

	buffer := bytes.NewBuffer([]byte{})
	if err = graph.Render(chart.PNG, buffer); err != nil {
		shared.Logger.Error(errors.Wrap(err, "failed to render chart"))
	}

	if img, _, err := image.Decode(buffer); err != nil {
		shared.Logger.Error(errors.Wrap(err, "failed to decode chart image"))
	} else {
		ctx.DrawImage(img, 0, int(float64(frameHeight/ 2)))
	}

	ShowFrame()
}

func RenderTextf(format string, a ...interface{}) {
	RenderText(fmt.Sprintf(format, a...))
}

// RenderTextWithIcon displays frame with `msg` text and `icon` image.
func RenderTextWithIcon(text, icon string) {
	initContext()

	var (
		iconImg image.Image
		iconPath = iconPath(icon)
		font, err = truetype.Parse(goregular.TTF)
		face = truetype.NewFace(font, &truetype.Options{
			Size: 14,
		})
	)

	if err != nil {
		shared.Logger.Error(errors.Wrap(err, "failed to create font face"))
		return
	}

	var (
		tw, th = ctx.MeasureString(text)
		tx, ty = float64(frameWidth/ 2) - tw / 2, float64(frameHeight/ 2) - th
	)

	if iconImg, err = gg.LoadPNG(iconPath); err != nil {
		shared.Logger.Error(errors.Wrapf(err, "failed to load from path '%s'", iconPath))
		return
	}

	var (
		ib = iconImg.Bounds()
		ix = int(math.Round(float64(frameWidth/ 2) - float64(ib.Dx()) / 2))
		iy = int(math.Round(float64(frameHeight / 2)))
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
		qrImg = qr.Image(frameHeight)
		x = int(math.Round(float64(frameWidth/ 2) - float64(qrImg.Bounds().Dx()) / 2))
		y = int(math.Round(float64(frameHeight/ 2) - float64(qrImg.Bounds().Dy()) / 2))
	)

	ctx.DrawImage(qrImg, x, y)
	ctx.Fill()

	ShowFrame()
}

// ShowFrame displays frame with rendered context.
func ShowFrame() {
	ctx.Fill()

	if !Available() {
		return
	}

	shared.MustExecute(func() error {
		return dev.DrawAndRefresh(ctx.Image())
	}, "failed to draw and refresh frame")
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
}

func clearFrame() {
	ctx.SetRGB(1, 1, 1)
	ctx.Clear()
	ctx.SetRGB(0, 0, 0)
}

func iconPath(icon string) string {
	return fmt.Sprintf("controllers/gui/assets/%s.png", icon)
}
