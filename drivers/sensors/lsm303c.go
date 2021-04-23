package sensors

import (
	"github.com/bskari/go-lsm303"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/peripherals"
	"github.com/timoth-y/iot-blockchain-sensorsys/model"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)


type (
	// LSM303Accelerometer defines accelerometer sensor device
	LSM303Accelerometer struct {
		conn *peripherals.I2C
		dev  *lsm303.Accelerometer
	}

	// LSM303Magnetometer defines magnetometer sensor device
	LSM303Magnetometer struct {
		conn *peripherals.I2C
		dev  *lsm303.Magnetometer
	}
)

func NewAccelerometerLSM303(addr uint16, bus int) *LSM303Accelerometer {
	return &LSM303Accelerometer{
		conn: peripherals.NewI2C(addr, bus),
	}
}

func NewMagnetometerLSM303(addr uint16, bus int) *LSM303Magnetometer {
	return &LSM303Magnetometer{
		conn: peripherals.NewI2C(addr, bus),
	}
}

func (s *LSM303Accelerometer) Init() (err error) {
	if err = s.conn.Init(); err != nil {
		shared.Logger.Error("connection init", err)
		return
	}

	if s.dev, err = lsm303.NewAccelerometer(s.conn.Bus,
		lsm303.WithAccelerometerSensorType(lsm303.LSM303C),
		lsm303.WithAccelerometerAddress(s.conn.Addr),
		lsm303.WithRange(lsm303.ACCELEROMETER_RANGE_2G),
	); err != nil {
		return
	}

	return
}

// ReadAxesG retrieves axes acceleration data as multiplications of G
func (s *LSM303Accelerometer) ReadAxesG() (model.Vector, error) {
	x, y, z, err := s.dev.SenseRaw(); if err != nil {
		return model.Vector{}, err
	}

	return model.Vector {
		X: round(float64(x) * scaleMultiplier, 4),
		Y: round(float64(y) * scaleMultiplier, 4),
		Z: round(float64(z) * scaleMultiplier, 4),
	}, nil
}

// ReadAxesMS2 parses data returned by GetAxesG and returns them in [m/s^2]
func (s *LSM303Accelerometer) ReadAxesMS2() (model.Vector, error) {
	x, y, z, err := s.dev.SenseRaw(); if err != nil {
		return model.Vector{}, err
	}

	return model.Vector {
		X: float64(x) * earthGravityMS2,
		Y: float64(y) * earthGravityMS2,
		Z: float64(z) * earthGravityMS2,
	}, nil
}

func (s *LSM303Accelerometer) ID() string {
	return "LSM303C-A"
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
	return s.conn.Active()
}

// Close disconnects from the device
func (s *LSM303Accelerometer) Close() error {
	return s.conn.Close()
}

func (s *LSM303Magnetometer) Init() (err error) {
	if err = s.conn.Init(); err != nil {
		return
	}

	if s.dev, err = lsm303.NewMagnetometer(s.conn.Bus,
		lsm303.WithMagnetometerSensorType(lsm303.LSM303C),
		lsm303.WithMagnetometerAddress(s.conn.Addr),
	); err != nil {
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

// ReadAxes parses data returned as magnetic force vector
func (s *LSM303Magnetometer) ReadTemperature() (float64, error) {
	t, err := s.dev.SenseRelativeTemperature(); if err != nil {
		return 0, err
	}

	return t.Celsius(), nil
}

func (s *LSM303Magnetometer) ID() string {
	return "LSM303C-M"
}

func (s *LSM303Magnetometer) Harvest(ctx *Context) {
	ctx.For(metrics.Magnetism).WriteWithError(toMagnitude(s.ReadAxes()))
	ctx.For(metrics.Temperature).WriteWithError(s.ReadTemperature())
}

func (s *LSM303Magnetometer) Metrics() []models.Metric {
	return []models.Metric{
		metrics.Magnetism,
		metrics.Temperature,
	}
}

func (s *LSM303Magnetometer) Verify() bool {
	return true
}

func (s *LSM303Magnetometer) Active() bool {
	return s.conn.Active()
}

// Close disconnects from the device
func (s *LSM303Magnetometer) Close() error {
	return s.conn.Close()
}
