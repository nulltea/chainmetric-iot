package peripherals

import (
	"github.com/pkg/errors"
	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
)

// SPI provides wrapper for SPI peripheral
type SPI struct {
	conn.Conn
	name string
	port spi.PortCloser
	active bool
}

func NewSPI(name string) *SPI {
	return &SPI{
		name: name,
	}
}

func (s *SPI) Init() (err error) {
	if err = s.InitPort(); err != nil {
		return err
	}

	if s.Conn, err = s.port.Connect(20 * physic.MegaHertz, spi.Mode0, 8); err != nil {
		return errors.Wrapf(err, "failed to connect vis SPI device on %s", s.name)
	}

	return
}

// InitPort initialises SPI port but not connects to it.
func (s *SPI) InitPort() (err error) {
	if s.port, err = spireg.Open(s.name); err != nil {
		return errors.Wrapf(err, "failed to open an SPI port on %s", s.name)
	}

	s.active = true

	return
}

func (s *SPI) Port() spi.Port {
	return s.port
}

func (s *SPI) Active() bool {
	return s.active
}

func (s *SPI) Close() error {
	s.active = false
	return s.port.Close()
}
