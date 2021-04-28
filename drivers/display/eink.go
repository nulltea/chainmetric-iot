package display

import (
	"image"

	"github.com/pkg/errors"
	"periph.io/x/periph/experimental/devices/epd"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripherals"
	"github.com/timoth-y/chainmetric-sensorsys/model/config"
)

// EInk is an implementation of Display driver for E-Ink 2.13" display.
type EInk struct {
	*epd.Dev
	*peripherals.SPI

	dc   *peripherals.GPIO
	cs   *peripherals.GPIO
	rst  *peripherals.GPIO
	busy *peripherals.GPIO

	config config.DisplayConfig
}

func NewEInk(config config.DisplayConfig) Display {
	return &EInk{
		SPI: peripherals.NewSPI(config.Bus),
		dc: peripherals.NewGPIO(config.DCPin),
		cs: peripherals.NewGPIO(config.CSPin),
		rst: peripherals.NewGPIO(config.ResetPin),
		busy:  peripherals.NewGPIO(config.BusyPin),
		config: config,
	}
}

func (d *EInk) Init() (err error) {
	if err = d.dc.Init(); err != nil {
		return errors.Wrap(err, "error during connecting to EInk display DC pin")
	}

	if err = d.cs.Init(); err != nil {
		return errors.Wrap(err, "error during connecting to EInk display CS pin")
	}

	if err = d.rst.Init(); err != nil {
		return errors.Wrap(err, "error during connecting to EInk display RST pin")
	}

	if err = d.busy.Init(); err != nil {
		return errors.Wrap(err, "error during connecting to EInk display BUSY pin")
	}

	if err = d.SPI.Init(); err != nil {
		return errors.Wrap(err, "error during connecting to EInk display via SPI")
	}

	d.Dev, err = epd.NewSPI(d.SPI.Port(), d.dc, d.cs, d.rst, d.busy, &epd.Opts{
		W: d.config.Width,
		H: d.config.Height,
		FullUpdate: epd.EPD2in13.FullUpdate,
		PartialUpdate: epd.EPD2in13.PartialUpdate,
	}); if err != nil {
		return errors.Wrap(err, "error during initialising to EInk display driver")
	}

	return
}

func (d *EInk) DrawImage(src image.Image) error {
	return d.Draw(d.Bounds(), src, image.Point{})
}
