package sensors

import (
	"github.com/bskari/go-lsm303"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/peripherals"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/sensor"
	"github.com/timoth-y/iot-blockchain-sensorsys/model"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)


type (
	// LSM303Accelerometer defines accelerometer sensor device
	LSM303Accelerometer struct {
		*peripherals.I2C
		dev  *lsm303.Accelerometer
	}

	// LSM303Magnetometer defines magnetometer sensor device
	LSM303Magnetometer struct {
		*peripherals.I2C
		dev  *lsm303.Magnetometer
	}
)

func NewAccelerometerLSM303(addr uint16, bus int) sensor.Sensor {
	return &LSM303Accelerometer{
		I2C: peripherals.NewI2C(addr, bus),
	}
}

func NewMagnetometerLSM303(addr uint16, bus int) sensor.Sensor {
	return &LSM303Magnetometer{
		I2C: peripherals.NewI2C(addr, bus),
	}
}

func (s *LSM303Accelerometer) Init() (err error) {
	if err = s.I2C.Init(); err != nil {
		shared.Logger.Error("connection init", err)
		return
	}

	if s.dev, err = lsm303.NewAccelerometer(s.Bus,
		lsm303.WithAccelerometerSensorType(lsm303.LSM303C),
		lsm303.WithAccelerometerAddress(s.Addr),
		lsm303.WithRange(lsm303.ACCELEROMETER_RANGE_2G),
	); err != nil {
		return
	}

	return
}

// ReadAxes retrieves axes acceleration data as multiplications of G
func (s *LSM303Accelerometer) ReadAxes() (model.Vector, error) {
	x, y, z, err := s.dev.SenseRaw(); if err != nil {
		return model.Vector{}, err
	}

	return model.Vector {
		X: round(float64(x) * scaleMultiplier, 4),
		Y: round(float64(y) * scaleMultiplier, 4),
		Z: round(float64(z) * scaleMultiplier, 4),
	}, nil
}

func (s *LSM303Accelerometer) ID() string {
	return "LSM303C-A"
}

func (s *LSM303Accelerometer) Harvest(ctx *sensor.Context) {
	ctx.For(metrics.Acceleration).WriteWithError(toMagnitude(s.ReadAxes()))
}

func (s *LSM303Accelerometer) Metrics() []models.Metric {
	return []models.Metric{
		metrics.Acceleration,
	}
}

func (s *LSM303Accelerometer) Verify() bool {
	return true
}

func (s *LSM303Magnetometer) Init() (err error) {
	if err = s.I2C.Init(); err != nil {
		return
	}

	if s.dev, err = lsm303.NewMagnetometer(s.Bus,
		lsm303.WithMagnetometerSensorType(lsm303.LSM303C),
		lsm303.WithMagnetometerAddress(s.Addr),
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

func (s *LSM303Magnetometer) Harvest(ctx *sensor.Context) {
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
