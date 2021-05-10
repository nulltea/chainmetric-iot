package display

// EInk display commands.
const (
	driverOutputControl            byte = 0x01
	boosterSoftStartControl        byte = 0x0C
	gateScanStartPosition          byte = 0x0F
	deepSleepMode                  byte = 0x10
	dataEntryModeSetting           byte = 0x11
	swReset                        byte = 0x12
	temperatureSensorControl       byte = 0x18
	masterActivation               byte = 0x20
	displayUpdateControl1          byte = 0x21
	displayUpdateControl2          byte = 0x22
	writeRAM                       byte = 0x24
	writeVcomRegister              byte = 0x2C
	writeLutRegister               byte = 0x32
	setDummyLinePeriod             byte = 0x3A
	setGateTime                    byte = 0x3B
	borderWaveformControl          byte = 0x3C
	setRAMXAddressStartEndPosition byte = 0x44
	setRAMYAddressStartEndPosition byte = 0x45
	autoWriteRamBW                 byte = 0x47
	setRAMXAddressCounter          byte = 0x4E
	setRAMYAddressCounter          byte = 0x4F
	terminateFrameReadWrite        byte = 0xFF
)

// LCD display commands.
const (
	NOP        = 0x00
	SWRESET    = 0x01
	RDDID      = 0x04
	RDDST      = 0x09
	SLPIN      = 0x10
	SLPOUT     = 0x11
	PTLON      = 0x12
	NORON      = 0x13
	INVOFF     = 0x20
	INVON      = 0x21
	DISPOFF    = 0x28
	DISPON     = 0x29
	CASET      = 0x2A
	RASET      = 0x2B
	RAMWR      = 0x2C
	RAMRD      = 0x2E
	PTLAR      = 0x30
	COLMOD     = 0x3A
	MADCTL     = 0x36
	MADCTL_MY  = 0x80
	MADCTL_MX  = 0x40
	MADCTL_MV  = 0x20
	MADCTL_ML  = 0x10
	MADCTL_RGB = 0x00
	MADCTL_BGR = 0x08
	MADCTL_MH  = 0x04
	RDID1      = 0xDA
	RDID2      = 0xDB
	RDID3      = 0xDC
	RDID4      = 0xDD
	FRMCTR1    = 0xB1
	RGBCTRL    = 0xB1
	FRMCTR2    = 0xB2
	PORCTRL    = 0xB2
	FRMCTR3    = 0xB3
	INVCTR     = 0xB4
	DISSET5    = 0xB6
	PWCTR1     = 0xC0
	PWCTR2     = 0xC1
	PWCTR3     = 0xC2
	PWCTR4     = 0xC3
	PWCTR5     = 0xC4
	VMCTR1     = 0xC5
	FRCTRL2    = 0xC6
	PWCTR6     = 0xFC
	GMCTRP1    = 0xE0
	GMCTRN1    = 0xE1
	GSCAN      = 0x45
	VSCRDEF    = 0x33
	VSCRSADD   = 0x37

	MAX_VSYNC_SCANLINES = 254
)
