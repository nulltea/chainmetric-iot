package gui

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-iot/shared"
	"golang.org/x/image/font/gofont/goregular"
)

var (
	batteryLevel = 100
)

func SetBatteryLevel(level int) {
	batteryLevel = level
}


func renderBatteryLevel() {
	var(
		icon = iconPath("battery")
		iconImg, err = gg.LoadPNG(iconPath("battery"))
		line = fmt.Sprintf("%d%%", batteryLevel)
	)

	if err != nil {
		shared.Logger.Error(errors.Wrapf(err, "failed to load from path '%s'", icon))
		return
	}

	var (
		ib = iconImg.Bounds()
		tx, th = ctx.MeasureString(line)
	)

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		shared.Logger.Error(errors.Wrap(err, "failed to parse font"))
	}

	ctx.SetFontFace(truetype.NewFace(font, &truetype.Options{
		Size: 13,
	}))

	ctx.DrawString(line, float64(frameWidth - ib.Dx()) - tx - 2.5, th)
	ctx.DrawImage(iconImg, int(frameWidth - ib.Dx()), 1)
}
