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
