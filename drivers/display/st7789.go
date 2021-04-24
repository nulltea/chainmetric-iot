package display

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"time"

	"github.com/pkg/errors"

	"github.com/timoth-y/iot-blockchain-sensorsys/model/config"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/peripherals"
)

// ST7789 is an open handle to the display controller.
type ST7789 struct {
	spi          *peripherals.SPI
	backlightPin *peripherals.GPIO
	resetPin     *peripherals.GPIO
	dcPin        *peripherals.GPIO

	rotation         Rotation
	frameRate        FrameRate
	width            int16
	height           int16
	batchLength      int32
	isBGR            bool
}

func NewST7789() *ST7789 {
	return &ST7789{}
}

func (d *ST7789) Init(config config.DisplayConfig) (err error) {
	d.dcPin = peripherals.NewGPIO(config.DCPin)
	if err := d.dcPin.Init(); err != nil {
		return errors.Wrap(err, "failed to connect DC pin")
	}

	d.backlightPin = peripherals.NewGPIO(config.BacklightPin)
	if err := d.backlightPin.Init(); err != nil {
		return errors.Wrap(err, "failed to connect backlight pin")
	}

	d.resetPin = peripherals.NewGPIO(config.ResetPin)
	if err := d.resetPin.Init(); err != nil {
		return errors.Wrap(err, "failed to connect reset pin")
	}

	if err = d.dcPin.Low(); err != nil {
		return err
	}

	d.spi = peripherals.NewSPI(config.Bus)

	if err = d.spi.Init(); err != nil {
		return errors.Wrap(err, "failed to connect via SPI bus")
	}

	if config.Width != 0 {
		d.width = int16(config.Width)
	} else {
		d.width = 240
	}

	if config.Height != 0 {
		d.height = int16(config.Height)
	} else {
		d.height = 240
	}

	if config.FrameRate != 0 {
		d.frameRate = FrameRate(config.FrameRate)
	} else {
		d.frameRate = FRAMERATE_60
	}

	if config.Rotation != 0 {
		d.rotation = Rotation(config.Rotation)
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

	return nil
}

// PowerOn the display
func (d *ST7789) PowerOn() error {
	return d.backlightPin.High()
}

// PowerOff the display
func (d *ST7789) PowerOff() error {
	return d.backlightPin.Low()
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

// FillRectangle fills a rectangle at a given coordinates with a color
func (d *ST7789) FillRectangle(x, y, width, height int16, c color.RGBA) error {
	k, i := d.Size()

	if x < 0 || y < 0 || width <= 0 || height <= 0 ||
		x >= k || (x+width) > k || y >= i || (y+height) > i {
		return errors.New("st7789: rectangle coordinates outside display area")
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
		return errors.New("st7789: rectangle coordinates outside display area")
	}

	if int32(width) * int32(height) != int32(len(buffer)) {
		return errors.New("st7789: buffer length does not match with rectangle size")
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

// FillScreen fills the screen with a given color
func (d *ST7789) FillScreen(c color.RGBA) {
	if d.rotation == NO_ROTATION || d.rotation == ROTATION_180 {
		d.FillRectangle(0, 0, d.width * 2, d.height * 2, c)
	} else {
		d.FillRectangle(0, 0, d.width * 2, d.height * 2, c)
	}
}

func (d *ST7789) DrawImage(img image.Image) {
	d.setWindow(0, 0, d.width, d.height)
	rect := d.Bounds()
	imgRect := img.Bounds()
	rgbaImg := image.NewRGBA(rect)
	centered := image.Point{
		X: (imgRect.Dx() - rect.Dx())/2,
		Y: (imgRect.Dy() - rect.Dx())/2,
	}
	draw.Draw(rgbaImg, rect, img, centered, draw.Src)

	buffer := make([]color.RGBA, 0)
	for i := 0; i < int(d.batchLength); i++ {
		for j := 0; j < int(d.batchLength); j++ {
			rgba := rgbaImg.RGBAAt(int(d.batchLength)-i, j)
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
		break
	case 1:
		madctl = MADCTL_MY | MADCTL_MV
		break
	case 2:
		break
	case 3:
		madctl = MADCTL_MX | MADCTL_MV
		break
	}
	if d.isBGR {
		madctl |= MADCTL_BGR
	}
	d.Command(MADCTL)
	d.Data(madctl)
}

// SetBGR changes the color mode (RGB/BGR)
func (d *ST7789) SetBGR(bgr bool) {
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
	return image.Rect(0, 0, int(d.width), int(d.height))
}

func (d *ST7789) SendData(c []byte) error {
	if err := d.dcPin.High(); err != nil {
		return err
	}
	return d.spi.Tx(c, nil)
}

func (d *ST7789) SendCommand(c []byte) error {
	if err := d.dcPin.Low(); err != nil {
		return err
	}
	return d.spi.Tx(c, nil)
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
	return d.spi.Close()
}

// RGBATo565 converts a color.RGBA to uint16 used in the display
func RGBATo565(c color.RGBA) uint16 {
	r, g, b, _ := c.RGBA()
	return uint16((r & 0xF800) +
		((g & 0xFC00) >> 5) +
		((b & 0xF800) >> 11))
}

func (d *ST7789) setWindow(x, y, w, h int16) {
	d.Command(CASET)
	d.SendData([]uint8{uint8(x >> 8), uint8(x), uint8((x + w - 1) >> 8), uint8(x + w - 1)})


	d.Command(RASET)
	d.SendData([]uint8{uint8(y >> 8), uint8(y), uint8((y + h - 1) >> 8), uint8(y + h - 1)})

	d.Command(RAMWR)
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
	d.setWindow(0, 0, d.width * 2, d.height * 2)
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
	d.resetPin.High()
	time.Sleep(50 * time.Millisecond)
	d.resetPin.Low()
	time.Sleep(50 * time.Millisecond)
	d.resetPin.High()
	time.Sleep(50 * time.Millisecond)
}
