package periphery

import (
	"github.com/pkg/errors"
	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
)

// SPI provides wrapper for SPI peripheral.
type SPI struct {
	conn.Conn
	name string
	port spi.PortCloser
	active bool
}

// NewSPI constructs new SPI driver instance.
func NewSPI(name string) *SPI {
	return &SPI{
		name: name,
	}
}

// Init performs SPI device initialization.
func (s *SPI) Init() (err error) {
	if s.port, err = spireg.Open(s.name); err != nil {
		return errors.Wrapf(err, "failed to open an SPI port on %s", s.name)
	}

	if s.Conn, err = s.port.Connect(20 * physic.MegaHertz, spi.Mode0, 8); err != nil {
		return errors.Wrapf(err, "failed to connect vis SPI device on %s", s.name)
	}

	s.active = true

	return
}

// SendCommandArgs sends `cmd` command with `data` as arguments on SPI device.
func (s *SPI) SendCommandArgs(cmd byte, data ...byte) error {
	if err := s.SendCommand(cmd); err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	return s.SendData(data...)
}

// SendCommand sends `cmd` command on SPI device.
func (s *SPI) SendCommand(cmd byte) error {
	if err := s.Tx([]byte{cmd}, nil); err != nil {
		return errors.Wrapf(err, "error during sending command 0x%X to SPI device", uint(cmd))
	}

	return nil
}

// SendData sends `data` on SPI device.
func (s *SPI) SendData(data ...byte) error {
	if len(data) == 0 {
		return nil
	}

	if err := s.Tx(data, nil); err != nil {
		return errors.Wrap(err, "error during sending data to SPI device")
	}

	return nil
}

// Port returns SPI port.
func (s *SPI) Port() spi.Port {
	return s.port
}

// Active checks whether the SPI device is connected and active.
func (s *SPI) Active() bool {
	return s.active
}

// Close closes connection to SPI device and clears allocated resources.
func (s *SPI) Close() error {
	s.active = false
	return s.port.Close()
}
