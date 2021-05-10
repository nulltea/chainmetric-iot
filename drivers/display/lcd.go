package display

import (
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"time"

	"github.com/pkg/errors"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/devices/ssd1306/image1bit"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripheries"
	"github.com/timoth-y/chainmetric-sensorsys/model/config"
)

// LCD is an implementation of Display driver for ST7789 2" display.
type LCD struct {
	*peripheries.SPI
	dc  *peripheries.GPIO
	cs  *peripheries.GPIO
	bl  *peripheries.GPIO
	rst *peripheries.GPIO

	rect image.Rectangle

	config config.DisplayConfig
}

// NewLCD creates new LCD driver instance by implementing Display interface.
func NewLCD(config config.DisplayConfig) Display {
	return &LCD{
		SPI:    peripheries.NewSPI(config.Bus),
		dc:     peripheries.NewGPIO(config.DCPin),
		cs:     peripheries.NewGPIO(config.CSPin),
		rst:    peripheries.NewGPIO(config.ResetPin),
		bl:     peripheries.NewGPIO(config.BusyPin),
		rect:   image.Rect(0, 0, config.Width, config.Height),
		config: config,
	}
}

func (d *LCD) Init() (err error) {
	if err = d.dc.Init(); err != nil {
		return errors.Wrap(err, "error during connecting to EInk display DC pin")
	}

	if err = d.cs.Init(); err != nil {
		return errors.Wrap(err, "error during connecting to EInk display CS pin")
	}

	if err = d.rst.Init(); err != nil {
		return errors.Wrap(err, "error during connecting to EInk display RST pin")
	}

	if err = d.bl.Init(); err != nil {
		return errors.Wrap(err, "error during connecting to EInk display BL pin")
	}

	if err = d.SPI.Init(); err != nil {
		return errors.Wrap(err, "error during connecting to LCD display via SPI")
	}

	if err := d.init(); err != nil {
		return errors.Wrap(err, "error during initialising to EInk display driver")
	}

	d.Reset()
	d.PowerOn()

	return nil
}

// PowerOn the display
func (d *LCD) PowerOn() error {
	return d.bl.High()
}

// Sleep the display
func (d *LCD) Sleep() error {
	return d.bl.Low()
}

// Draw implements display.Drawer.
func (d *LCD) Draw(r image.Rectangle, src image.Image, sp image.Point) error {
	var (
		xStart = sp.X
		yStart = sp.Y
		imageW = r.Dx() & 0xF8
		imageH = r.Dy()
		w      = d.rect.Dx()
		h      = d.rect.Dy()
	)


	xEnd := xStart + imageW - 1
	if xStart+imageW >= w {
		xEnd = w - 1
	}

	yEnd := yStart + imageH - 1
	if yStart+imageH >= h {
		yEnd = h - 1
	}

	if err := d.setWindow(xStart, yStart, xEnd, yEnd); err != nil {
		return err
	}

	next := image1bit.NewVerticalLSB(d.rect)
	draw.Src.Draw(next, r, src, sp)
	var byteToSend byte = 0x00
	for y := yStart; y < yEnd+1; y++ {
		if err := d.SendCommandArgs(writeRAM); err != nil {
			return err
		}
		for x := xStart; x < xEnd+1; x++ {
			bit := next.BitAt(x-xStart, y-yStart)

			if bit {
				byteToSend |= 0x80 >> (uint32(x) % 8)
			}

			if x%8 == 7 {
				if err := d.SendData(byteToSend); err != nil {
					return err
				}
				byteToSend = 0x00
			}
		}
	}

	return nil
}

// DrawImage sends `src` image binary representation to LCD display buffer.
// Use Refresh() or DrawAndRefresh() to display image.
func (d *LCD) DrawImage(src image.Image) error {
	return d.Draw(d.Bounds(), src, image.Point{})
}


func (d *LCD) Clear() error {
	return d.ResetFrame(0xFF)
}

// Bounds implements display.Drawer. Min is guaranteed to be {0, 0}.
func (d *LCD) Bounds() image.Rectangle {
	return d.rect
}

// SendCommandArgs overrides peripheries.SPI send command with args method
// by additionally sending signals to DC GPIO pins.
func (d *LCD) SendCommandArgs(cmd byte, data ...byte) error {
	if !d.Active() {
		return nil
	}

	if err := d.SendCommand(cmd); err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	return d.SendData(data...)
}

// SendCommand overrides peripheries.SPI send command method
// by additionally sending signals to DC GPIO pins.
func (d *LCD) SendCommand(cmd byte) (err error) {
	if !d.Active() {
		return
	}

	if err := d.dc.Out(gpio.Low); err != nil {
		return errors.Wrapf(err, "error during sending %s signal to %s", d.dc, gpio.Low)
	}

	return d.SPI.SendCommand(cmd)
}

// SendData overrides peripheries.SPI send data method
// by additionally sending signals to DC GPIO pins.
func (d *LCD) SendData(data ...byte) (err error) {
	if !d.Active() {
		return nil
	}

	if len(data) == 0 {
		return nil
	}

	if err := d.dc.Out(gpio.High); err != nil {
		return errors.Wrapf(err, "error during sending %s signal to %s", d.dc, gpio.High)
	}

	return d.SPI.SendData(data...)
}

func (d *LCD) Active() bool {
	return d.SPI.Active()
}

func (d *LCD) Close() error {
	d.Sleep()
	return d.SPI.Close()
}

// ResetFrame clear the frame memory with the specified color.
// this won't update the display.
func (d *LCD) ResetFrame(color byte) error {
	var (
		w = d.rect.Dx()
		h = d.rect.Dy()
	)

	// send the color data
	for i := 0; i < (w / 8 * h); i++ {
		if err := d.SendData(color); err != nil {
			return err
		}
	}

	return nil
}


func (d *LCD) setWindow(x, y, w, h int) error {
	if err := d.SendCommandArgs(CASET, uint8(x >> 8), uint8(x), uint8((x + w - 1) >> 8), uint8(x + w - 1)); err != nil {
		return err
	}

	if err := d.SendCommandArgs(RASET, uint8(y >> 8), uint8(y), uint8((y + h - 1) >> 8), uint8(y + h - 1)); err != nil {
		return err
	}

	if err := d.SendCommand(RAMWR); err != nil {
		return err
	}

	return nil
}

func (d *LCD) init() error {
	// Common initialization
	d.SendCommand(SWRESET)
	time.Sleep(150 * time.Millisecond)

	// Exit sleep mode
	d.SendCommand(SLPOUT)
	time.Sleep(500 * time.Millisecond)

	// Set color mode: 16-bit color
	d.SendCommandArgs(COLMOD, 0x55)
	time.Sleep(10 * time.Millisecond)

	// Clear screen
	d.setWindow(0, 0, d.config.Width, d.config.Height)
	d.Clear()

	// Ready to display
	d.SendCommand(INVON)
	time.Sleep(10 * time.Millisecond)

	d.SendCommand(NORON)
	time.Sleep(10 * time.Millisecond)

	d.SendCommand(DISPON)
	time.Sleep(10 * time.Millisecond)

	return nil
}

func (d *LCD) Reset() error {
	if err := d.rst.High(); err != nil {
		return err
	}
	time.Sleep(50 * time.Millisecond)

	if err := d.rst.Low(); err != nil {
		return err
	}
	time.Sleep(50 * time.Millisecond)

	if err := d.rst.High(); err != nil {
		return err
	}
	time.Sleep(50 * time.Millisecond)

	return nil
}
