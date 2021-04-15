package sensors

import (
	"github.com/bskari/go-lsm303"
	"github.com/timoth-y/iot-blockchain-contracts/models"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"

	"github.com/timoth-y/iot-blockchain-sensorsys/model"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)


type (
	// LSM303Accelerometer defines accelerometer sensor device
	LSM303Accelerometer struct {
		addr uint8
		n    int
		bus  i2c.BusCloser
		dev  *lsm303.Accelerometer
	}

	// LSM303Magnetometer defines magnetometer sensor device
	LSM303Magnetometer struct {
		addr uint8
		n    int
		bus  i2c.BusCloser
		dev  *lsm303.Magnetometer
	}
)

func NewAccelerometerLSM303(addr uint8, bus int) *LSM303Accelerometer {
	return &LSM303Accelerometer{
		addr: addr,
		n:    bus,
	}
}

func NewMagnetometerLSM303(addr uint8, bus int) *LSM303Magnetometer {
	return &LSM303Magnetometer{
		addr: addr,
		n:    bus,
	}
}

func (s *LSM303Accelerometer) Init() (err error) {
	s.bus, err = i2creg.Open(shared.NtoI2cBusName(s.n)); if err != nil {
		return
	}

	s.dev, err = lsm303.NewAccelerometer(s.bus, &lsm303.DefaultAccelerometerOpts); if err != nil {
		return
	}

	return
}

// ReadAxesG retrieves axes acceleration data as multiplications of G
func (s *LSM303Accelerometer) ReadAxesG() (model.Vector, error) {
	x, y, z, err := s.dev.Sense(); if err != nil {
		return model.Vector{}, err
	}

	return model.Vector {
		X: round(float64(x) * earthGravityMS2, 4),
		Y: round(float64(y) * earthGravityMS2, 4),
		Z: round(float64(z) * earthGravityMS2, 4),
	}, nil
}

// ReadAxesMS2 parses data returned by GetAxesG and returns them in [m/s^2]
func (s *LSM303Accelerometer) ReadAxesMS2() (model.Vector, error) {
	x, y, z, err := s.dev.Sense(); if err != nil {
		return model.Vector{}, err
	}

	return model.Vector {
		X: float64(x),
		Y: float64(y),
		Z: float64(z),
	}, nil
}

func (s *LSM303Accelerometer) ID() string {
	return "LSM303"
}

func (s *LSM303Accelerometer) Harvest(ctx *Context) {
	ctx.For(metrics.AccelerationInG).WriteWithError(toMagnitude(s.ReadAxesG()))
	ctx.For(metrics.AccelerationInMS2).WriteWithError(toMagnitude(s.ReadAxesMS2()))
}

func (s *LSM303Accelerometer) Metrics() []models.Metric {
	return []models.Metric{
		metrics.AccelerationInG,
		metrics.AccelerationInMS2,
	}
}

func (s *LSM303Accelerometer) Verify() bool {
	return true
}

func (s *LSM303Accelerometer) Active() bool {
	return s.dev != nil
}

// Close disconnects from the device
func (s *LSM303Accelerometer) Close() error {
	defer s.clean()
	return s.bus.Close()
}

func (s *LSM303Accelerometer) clean() {
	s.dev = nil
}

func (s *LSM303Magnetometer) Init() (err error) {
	s.bus, err = i2creg.Open(shared.NtoI2cBusName(s.n)); if err != nil {
		return
	}

	s.dev, err = lsm303.NewMagnetometer(s.bus, &lsm303.DefaultMagnetometerOpts); if err != nil {
		return
	}

	return
}

// ReadAxes parses data returned as magnetic force vector
func (s *LSM303Magnetometer) ReadAxes() (model.Vector, error) {
	x, y, z, err := s.dev.SenseRaw(); if err != nil {
		return model.Vector{}, err
	}

	return model.Vector {
		X: float64(x),
		Y: float64(y),
		Z: float64(z),
	}, nil
}

func (s *LSM303Magnetometer) ID() string {
	return "LSM303"
}

func (s *LSM303Magnetometer) Harvest(ctx *Context) {
	ctx.For(metrics.Magnetism).WriteWithError(toMagnitude(s.ReadAxes()))
}

func (s *LSM303Magnetometer) Metrics() []models.Metric {
	return []models.Metric{
		metrics.AccelerationInG,
		metrics.AccelerationInMS2,
	}
}

func (s *LSM303Magnetometer) Verify() bool {
	return true
}

func (s *LSM303Magnetometer) Active() bool {
	return s.dev != nil
}

// Close disconnects from the device
func (s *LSM303Magnetometer) Close() error {
	defer func() {
		s.dev = nil
	}()

	return s.bus.Close()
}
