package sensors

import (
	"fmt"

	"github.com/cgxeiji/max3010x"
)

type MAX30102 struct {
	addr uint8
	bus int
	dev *max3010x.Device
}

func NewMAX30102(addr uint8, bus int) *MAX30102 {
	return &MAX30102{
		addr: addr,
		bus: bus,
	}
}

func (s *MAX30102) Init() (err error) {
	s.dev, err = max3010x.NewOnBus(fmt.Sprintf("/dev/i2c-%d", s.bus)); if err != nil {
		return
	}

	if err = s.dev.Startup(); err != nil {
		return err
	}

	return
}

func (s *MAX30102) Read() (bpm float64, o2 float64, err error) {
	if bpm, err = s.ReadHeartRate(); err != nil {
		return 0, 0, err
	}
	if o2, err = s.ReadSpO2(); err != nil {
		return bpm, 0, err
	}

	return
}

func (s *MAX30102) ReadHeartRate() (float64, error) {
	return s.dev.HeartRate()
}

func (s *MAX30102) ReadSpO2() (float64, error) {
	return s.dev.SpO2()
}

func (s *MAX30102) Close() error {
	s.dev.Close()
	return nil
}



