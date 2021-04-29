package display

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"time"

	"github.com/pkg/errors"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/devices/ssd1306/image1bit"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripherals"
	"github.com/timoth-y/chainmetric-sensorsys/model/config"
)

// EInk is an implementation of Display driver for E-Ink 2.13" display.
type EInk struct {
	*peripherals.SPI

	dc   *peripherals.GPIO
	cs   *peripherals.GPIO
	rst  *peripherals.GPIO
	busy *peripherals.GPIO

	rect image.Rectangle

	config config.DisplayConfig
}

func NewEInk(config config.DisplayConfig) Display {
	return &EInk{
		SPI: peripherals.NewSPI(config.Bus),
		dc: peripherals.NewGPIO(config.DCPin),
		cs: peripherals.NewGPIO(config.CSPin),
		rst: peripherals.NewGPIO(config.ResetPin),
		busy:  peripherals.NewGPIO(config.BusyPin),
		rect: image.Rect(0, 0, config.Width, config.Height),
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

	if err := d.init(); err != nil {
		return errors.Wrap(err, "error during initialising to EInk display driver")
	}

	d.Clear()

	return
}

// DrawRaw implements display.Drawer.
func (d *EInk) DrawRaw(r image.Rectangle, src image.Image, sp image.Point) error {
	xStart := sp.X
	yStart := sp.Y
	imageW := r.Dx() & 0xF8
	imageH := r.Dy()
	w := d.rect.Dx()
	h := d.rect.Dy()

	xEnd := xStart + imageW - 1
	if xStart+imageW >= w {
		xEnd = w - 1
	}

	yEnd := yStart + imageH - 1
	if yStart+imageH >= h {
		yEnd = h - 1
	}

	if err := d.setMemoryArea(xStart, yStart, xEnd, yEnd); err != nil {
		return err
	}

	next := image1bit.NewVerticalLSB(d.rect)
	draw.Src.Draw(next, r, src, sp)
	var byteToSend byte = 0x00
	for y := yStart; y < yEnd+1; y++ {
		if err := d.setMemoryPointer(xStart, y); err != nil {
			return err
		}
		if err := d.sendCommand(writeRAM); err != nil {
			return err
		}
		for x := xStart; x < xEnd+1; x++ {
			bit := next.BitAt(x-xStart, y-yStart)
			if bit {
				byteToSend |= 0x80 >> (uint32(x) % 8)
			}
			if x%8 == 7 {
				if err := d.sendData(byteToSend); err != nil {
					return err
				}
				byteToSend = 0x00
			}
		}
	}

	return nil
}

func (d *EInk) Draw(src image.Image) error {
	return d.DrawRaw(src.Bounds(), src, image.Point{0, 0})
}

func (d *EInk) DrawAndRefresh(src image.Image) error {
	if err := d.Draw(src); err != nil {
		return err
	}

	return d.Refresh()
}

// ResetFrameMemory clear the frame memory with the specified color.
// this won't update the display.
func (d *EInk) ResetFrameMemory(color byte) error {
	w := d.rect.Dx()
	h := d.rect.Dy()
	if err := d.setMemoryArea(0, 0, w-1, h-1); err != nil {
		return err
	}
	if err := d.setMemoryPointer(0, 0); err != nil {
		return err
	}
	if err := d.sendCommand(writeRAM); err != nil {
		return err
	}

	// send the color data
	for i := 0; i < (w / 8 * h); i++ {
		if err := d.sendData(color); err != nil {
			return err
		}
	}

	return nil
}

// Refresh updates the display.
//
// There are 2 memory areas embedded in the e-paper display but once
// this function is called, the next action of SetFrameMemory or ClearFrame
// will set the other memory area.
func (d *EInk) Refresh() error {
	if err := d.sendCommand(displayUpdateControl2, 0xC4); err != nil {
		return err
	}


	if err := d.sendCommand(masterActivation); err != nil {
		return err
	}

	if err := d.sendCommand(terminateFrameReadWrite); err != nil {
		return err
	}

	d.waitUntilIdle()
	return nil
}

// Clear clears the display.
func (d *EInk) Clear() error {
	return d.ResetFrameMemory(0xFF)
}

// Sleep after this command is transmitted, the chip would enter the
// deep-sleep mode to save power.
//
// The deep sleep mode would return to standby by hardware reset.
// You can use Reset() to awaken and Init to re-initialize the device.
func (d *EInk) Sleep() error {
	if err := d.sendCommand(deepSleepMode); err != nil {
		return err
	}

	d.waitUntilIdle()
	return nil
}

// Reset can be also used to awaken the device
func (d *EInk) Reset() (err error) {
	if err = d.rst.Out(gpio.High); err != nil {
		return
	}
	time.Sleep(200 * time.Millisecond)

	if err = d.rst.Out(gpio.Low); err != nil {
		return
	}
	time.Sleep(200 * time.Millisecond)

	if err = d.rst.Out(gpio.High); err != nil {
		return
	}
	time.Sleep(200 * time.Millisecond)

	return
}

// ColorModel implements display.Drawer.
// It is a one bit color model, as implemented by image1bit.Bit.
func (d *EInk) ColorModel() color.Model {
	return image1bit.BitModel
}

// Bounds implements display.Drawer. Min is guaranteed to be {0, 0}.
func (d *EInk) Bounds() image.Rectangle {
	return d.rect
}

func (d *EInk) init() error {
	if err := d.Reset(); err != nil {
		return err
	}

	if err := d.sendCommand(swReset); err != nil {
		return err
	}

	d.waitUntilIdle()

	d.sendCommand(autoWriteRamBW, 0xF7)
	d.waitUntilIdle()

	if err := d.sendCommand(driverOutputControl,
		byte((d.config.Width - 1) & 0xFF),
		byte(((d.config.Height - 1) >> 8) & 0xFF), 0x01); err != nil {
		return err
	}

	if err := d.sendCommand(boosterSoftStartControl, 0xAE, 0xC7, 0xC3, 0xC0, 0x40); err != nil {
		return err
	}

	if err := d.sendCommand(dataEntryModeSetting, 0x01); err != nil {
		return err
	}

	if err := d.setMemoryArea(0, 0, d.config.Width - 1, d.config.Height - 1); err != nil {
		return err
	}

	if err := d.sendCommand(borderWaveformControl, 0x01); err != nil {
		return err
	}

	if err := d.sendCommand(temperatureSensorControl, 0x80); err != nil {
		return err
	}

	if err := d.sendCommand(displayUpdateControl2, 0xB1); err != nil {
		return err
	}

	if err := d.sendCommand(masterActivation); err != nil {
		return err
	}

	d.waitUntilIdle()

	if err := d.sendCommand(displayUpdateControl2, 0xB1); err != nil {
		return err
	}

	if err := d.sendCommand(masterActivation); err != nil {
		return err
	}

	return d.setMemoryPointer(0, 0)
}

func (d *EInk) setMemoryPointer(x, y int) error {
	if err := d.sendCommand(setRAMXAddressCounter, byte((x >> 3) & 0xFF)); err != nil {
		return err
	}

	if err := d.sendCommand(setRAMYAddressCounter, byte(y & 0xFF), byte((y >> 8) & 0xFF)); err != nil {
		return err
	}

	d.waitUntilIdle()

	return nil
}

func (d *EInk) waitUntilIdle() {
	for d.busy.Read() == gpio.High {
		time.Sleep(100 * time.Millisecond)
	}
}

func (d *EInk) setMemoryArea(xStart, yStart, xEnd, yEnd int) error {
	if err := d.sendCommand(setRAMXAddressStartEndPosition,
		byte((xStart >> 3) & 0xFF),
		byte((xEnd >> 3) & 0xFF),
	); err != nil {
		return err
	}

	return d.sendCommand(setRAMYAddressStartEndPosition,
		byte(yStart & 0xFF),
		byte((yStart >> 8) & 0xFF),
		byte(yEnd & 0xFF),
		byte((yEnd >> 8) & 0xFF),
	)
}


func (d *EInk) sendCommand(cmd byte, data ...byte) error {
	if err := d.writeCommand(cmd); err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	return d.sendData(data...)
}


func (d *EInk) writeCommand(p byte) (err error) {
	if err := d.dc.Out(gpio.Low); err != nil {
		return fmt.Errorf("%v.Out(%v) = %w", d.dc.String(), gpio.Low.String(), err)
	}

	if err := d.cs.Out(gpio.Low); err != nil {
		return fmt.Errorf("%v.Out(%v) = %w", d.cs.String(), gpio.Low.String(), err)
	}

	defer func() {
		if err2 := d.cs.Out(gpio.High); err2 != nil {
			err = fmt.Errorf("%v.Out(%v) = %w, already had error %v", d.cs.String(), gpio.High, err2, err)
		}
	}()

	if err := d.Tx([]byte{p}, nil); err != nil {
		return fmt.Errorf("sending command 0x%X: %w", uint16(p), err)
	}

	return nil
}

func (d *EInk) sendData(p ...byte) (err error) {
	if len(p) == 0 {
		return nil
	}
	if err := d.cs.Out(gpio.Low); err != nil {
		return fmt.Errorf("%v.Out(%v) = %w", d.cs.String(), gpio.Low.String(), err)
	}
	if err := d.dc.Out(gpio.High); err != nil {
		return fmt.Errorf("%v.Out(%v) = %w", d.dc.String(), gpio.High.String(), err)
	}
	defer func() {
		if e := d.cs.Out(gpio.High); e != nil {
			err = fmt.Errorf("already had err %q, and got e: %w", err, e)
		}
	}()

	if err := d.Tx(p, nil); err != nil {
		return  err
	}
	return nil
}

