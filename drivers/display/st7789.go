package display

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"time"

	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"

	"github.com/timoth-y/iot-blockchain-sensorsys/config"
	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

// ST7789 is an open handle to the display controller.
type ST7789 struct {
	bus    string
	pin    string
	port   spi.PortCloser
	conn   conn.Conn
	dc     gpio.PinOut
	rect   image.Rectangle

	rotation                      Rotation
	frameRate                     FrameRate
	width                         int16
	height                        int16
	rowOffsetCfg, rowOffset       int16
	columnOffset, columnOffsetCfg int16
	isBGR                         bool
	batchLength                   int32
	backlightPin                  gpio.PinIO
	resetPin                      gpio.PinIO
}

func NewST7789(bus string, pin string) *ST7789 {
	return &ST7789{
		bus:          bus,
		pin:          pin,
		backlightPin: gpioreg.ByName("GPIO18"),
		resetPin:     gpioreg.ByName("GPIO15"),
	}
}

func (d *ST7789) Init(cnf config.DisplayConfig) (err error) {
	if d.port, err = spireg.Open(d.bus); err != nil {
		return err
	}

	d.dc = gpioreg.ByName(d.pin); if d.dc == gpio.INVALID {
		return errors.New("ssd1306: use nil for dc to use 3-wire mode, do not use gpio.INVALID")
	}

	if err = d.dc.Out(gpio.Low); err != nil {
		return err
	}

	if d.conn, err = d.port.Connect(80*physic.MegaHertz, spi.Mode0, 8); err != nil {
		return err
	}

	if cnf.Width != 0 {
		d.width = int16(cnf.Width)
	} else {
		d.width = 240
	}

	if cnf.Height != 0 {
		d.height = int16(cnf.Height)
	} else {
		d.height = 240
	}

	if cnf.FrameRate != 0 {
		d.frameRate = FrameRate(cnf.FrameRate)
	} else {
		d.frameRate = FRAMERATE_60
	}

	if cnf.Rotation != 0 {
		d.rotation = Rotation(cnf.Rotation)
	} else {
		d.rotation = NO_ROTATION
	}

	d.batchLength = int32(d.width)
	if d.height > d.width {
		d.batchLength = int32(d.height)
	}
	d.batchLength += d.batchLength & 1

	d.reset()
	d.setup()

	//d.SetRotation(d.rotation)

	return nil
}

// PowerOn the display
func (d *ST7789) PowerOn() error {
	return d.backlightPin.Out(gpio.High)
}

// PowerOff the display
func (d *ST7789) PowerOff() error {
	return d.backlightPin.Out(gpio.Low)
}

func (d *ST7789) setWindow(x, y, w, h int16) {
	x += d.columnOffset
	y += d.rowOffset

	d.Command(CASET)
	d.SendData([]uint8{uint8(x >> 8), uint8(x), uint8((x + w - 1) >> 8), uint8(x + w - 1)})

	d.Command(RASET)
	d.SendData([]uint8{uint8(y >> 8), uint8(y), uint8((y + h - 1) >> 8), uint8(y + h - 1)})

	d.Command(RAMWR)
}

// FillRectangle fills a rectangle at a given coordinates with a color
func (d *ST7789) FillRectangle(x, y, width, height int16, c color.RGBA) error {
	k, i := d.Size()
	if x < 0 || y < 0 || width <= 0 || height <= 0 ||
		x >= k || (x+width) > k || y >= i || (y+height) > i {
		return errors.New("rectangle coordinates outside display area")
	}
	d.setWindow(x, y, width, height)
	c565 := RGBATo565(c)
	c1 := uint8(c565 >> 8)
	c2 := uint8(c565)

	data := make([]uint8, d.batchLength * 2)
	for i := int32(0); i < d.batchLength; i++ {
		data[i*2] = c1
		data[i*2+1] = c2
	}
	j := int32(width) * int32(height)
	for j > 0 {
		if j >= d.batchLength {
			d.SendData(data)
		} else {
			d.SendData(data[:j*2])
		}
		j -= d.batchLength
	}
	return nil
}

// FillRectangleWithBuffer fills buffer with a rectangle at a given coordinates.
func (d *ST7789) FillRectangleWithBuffer(x, y, width, height int16, buffer []color.RGBA) error {
	i, j := d.Size()
	if x < 0 || y < 0 || width <= 0 || height <= 0 ||
		x >= i || (x+width) > i || y >= j || (y+height) > j {
		return errors.New("rectangle coordinates outside display area")
	}
	if int32(width)*int32(height) != int32(len(buffer)) {
		return errors.New("buffer length does not match with rectangle size")
	}
	d.setWindow(x, y, width, height)

	k := int32(width) * int32(height)
	data := make([]uint8, d.batchLength*2)
	offset := int32(0)
	for k > 0 {
		for i := int32(0); i < d.batchLength; i++ {
			if offset+i < int32(len(buffer)) {
				c565 := RGBATo565(buffer[offset+i])
				c1 := uint8(c565 >> 8)
				c2 := uint8(c565)
				data[i*2] = c1
				data[i*2+1] = c2
			}
		}
		if k >= d.batchLength {
			d.SendData(data)
		} else {
			d.SendData(data[:k*2])
		}
		k -= d.batchLength
		offset += d.batchLength
	}
	return nil
}

// SetPixel sets a pixel in the screen
func (d *ST7789) SetPixel(x int16, y int16, c color.RGBA) {
	if x < 0 || y < 0 ||
		(((d.rotation == NO_ROTATION || d.rotation == ROTATION_180) && (x >= d.width || y >= d.height)) ||
			((d.rotation == ROTATION_90 || d.rotation == ROTATION_270) && (x >= d.height || y >= d.width))) {
		return
	}
	d.FillRectangle(x, y, 1, 1, c)
}

// FillScreen fills the screen with a given color
func (d *ST7789) FillScreen(c color.RGBA) {
	if d.rotation == NO_ROTATION || d.rotation == ROTATION_180 {
		d.FillRectangle(0, 0, d.width, d.height, c)
	} else {
		d.FillRectangle(0, 0, d.height, d.width, c)
	}
}

func (d *ST7789) DrawRAW(reader io.Reader) {
	d.setWindow(0, 0, d.width, d.height)
	img, _, err := image.Decode(reader)
	if err != nil {
		shared.Logger.Error(err)
		return
	}

	d.DrawImage(img)
}

func (d *ST7789) DrawImage(img image.Image) {
	d.setWindow(0, 0, d.width, d.height)
	rect := img.Bounds()
	rgbaImg := image.NewRGBA(rect)
	draw.Draw(rgbaImg, rect, img, rect.Min, draw.Src)


	buffer := make([]color.RGBA, 0)
	for i := 0; i < int(d.batchLength); i++ {
		for j := 0; j < int(d.batchLength); j++ {
			rgba := rgbaImg.RGBAAt(int(d.width)-i, j)
			buffer = append(buffer, rgba)
		}
	}

	d.FillRectangleWithBuffer(0, 0, d.width, d.height, buffer)
}

// DrawFastVLine draws a vertical line faster than using SetPixel
func (d *ST7789) DrawFastVLine(x, y0, y1 int16, c color.RGBA) {
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	d.FillRectangle(x, y0, 1, y1-y0+1, c)
}

// DrawFastHLine draws a horizontal line faster than using SetPixel
func (d *ST7789) DrawFastHLine(x0, x1, y int16, c color.RGBA) {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	d.FillRectangle(x0, y, x1-x0+1, 1, c)
}

// SetRotation changes the rotation of the device (clock-wise)
func (d *ST7789) SetRotation(rotation Rotation) {
	madctl := uint8(0)
	switch rotation % 4 {
	case 0:
		madctl = MADCTL_MX | MADCTL_MY
		d.rowOffset = d.rowOffsetCfg
		d.columnOffset = d.columnOffsetCfg
		break
	case 1:
		madctl = MADCTL_MY | MADCTL_MV
		d.rowOffset = d.columnOffsetCfg
		d.columnOffset = d.rowOffsetCfg
		break
	case 2:
		d.rowOffset = 0
		d.columnOffset = 0
		break
	case 3:
		madctl = MADCTL_MX | MADCTL_MV
		d.rowOffset = 0
		d.columnOffset = 0
		break
	}
	if d.isBGR {
		madctl |= MADCTL_BGR
	}
	d.Command(MADCTL)
	d.Data(madctl)
}

// IsBGR changes the color mode (RGB/BGR)
func (d *ST7789) IsBGR(bgr bool) {
	d.isBGR = bgr
}


// Size returns the current size of the display.
func (d *ST7789) Size() (w, h int16) {
	if d.rotation == NO_ROTATION || d.rotation == ROTATION_180 {
		return d.width, d.height
	}
	return 240, 240
}

// Bounds implements display.Drawer. Min is guaranteed to be {0, 0}.
func (d *ST7789) Bounds() image.Rectangle {
	return d.rect
}

// Invert the display (black on white vs white on black).
func (d *ST7789) Invert(blackOnWhite bool) {
	b := byte(0xA6)
	if blackOnWhite {
		b = 0xA7
	}
	d.Command(b)
}

// RGBATo565 converts a color.RGBA to uint16 used in the display
func RGBATo565(c color.RGBA) uint16 {
	r, g, b, _ := c.RGBA()
	return uint16((r & 0xF800) +
		((g & 0xFC00) >> 5) +
		((b & 0xF800) >> 11))
}


func (d *ST7789) SendData(c []byte) error {
	if err := d.dc.Out(gpio.High); err != nil {
		return err
	}
	return d.conn.Tx(c, nil)
}

func (d *ST7789) SendCommand(c []byte) error {
	if err := d.dc.Out(gpio.Low); err != nil {
		return err
	}
	return d.conn.Tx(c, nil)
}


// Command sends a command to the device
func (d *ST7789) Command(cmd uint8) {
	d.SendCommand([]byte{cmd})
}

// Data sends data to the device
func (d *ST7789) Data(data uint8) {
	d.SendData([]byte{data})
}

func (d *ST7789) Close() error {
	d.PowerOff()
	return d.port.Close()
}

func (d *ST7789) setup() {
	// Common initialization
	d.Command(SWRESET)
	time.Sleep(150 * time.Millisecond)

	// Exit sleep mode
	d.Command(SLPOUT)
	time.Sleep(500 * time.Millisecond)

	// Set color mode: 16-bit color
	d.Command(COLMOD)
	d.Data(0x55)
	time.Sleep(10 * time.Millisecond)

	// Set orientation
	d.SetRotation(d.rotation)

	// Clear screen
	d.setWindow(0, 0, d.width, d.height)
	d.FillScreen(color.RGBA{A: 255})

	// Frame rate for normal mode: 60Hz
	d.Command(FRCTRL2)
	d.Data(uint8(d.frameRate))

	// Ready to display
	d.Command(INVON)
	time.Sleep(10 * time.Millisecond)

	d.Command(NORON)
	time.Sleep(10 * time.Millisecond)

	d.Command(DISPON)
	time.Sleep(10 * time.Millisecond)
}

func (d *ST7789) reset() {
	d.resetPin.Out(gpio.High)
	time.Sleep(50 * time.Millisecond)
	d.resetPin.Out(gpio.Low)
	time.Sleep(50 * time.Millisecond)
	d.resetPin.Out(gpio.High)
	time.Sleep(50 * time.Millisecond)
}
