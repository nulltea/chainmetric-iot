package sensors

import (
	"sync"

	"github.com/bskari/go-lsm303"
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-iot/core/dev/sensor"
	"github.com/timoth-y/chainmetric-iot/drivers/periphery"
	"github.com/timoth-y/chainmetric-iot/model"
)


var (
	lsm303cAccelerometerMutex = &sync.Mutex{}
	lsm303cMagnetometerMutex  = &sync.Mutex{}
)

type (
	// LSM303Accelerometer defines accelerometer sensor device
	LSM303Accelerometer struct {
		*periphery.I2C
		dev  *lsm303.Accelerometer
	}

	// LSM303Magnetometer defines magnetometer sensor device
	LSM303Magnetometer struct {
		*periphery.I2C
		dev  *lsm303.Magnetometer
	}
)

func NewAccelerometerLSM303(addr uint16, bus int) sensor.Sensor {
	return &LSM303Accelerometer{
		I2C: periphery.NewI2C(addr, bus, periphery.WithMutex(lsm303cAccelerometerMutex)),
	}
}

func NewMagnetometerLSM303(addr uint16, bus int) sensor.Sensor {
	return &LSM303Magnetometer{
		I2C: periphery.NewI2C(addr, bus, periphery.WithMutex(lsm303cMagnetometerMutex)),
	}
}

func (s *LSM303Accelerometer) Init() (err error) {
	if err = s.I2C.Init(); err != nil {
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
	s.Lock()
	defer s.Unlock()

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
	ctx.WriterFor(metrics.Acceleration).WriteWithError(toMagnitude(s.ReadAxes()))
}

func (s *LSM303Accelerometer) Metrics() []models.Metric {
	return []models.Metric{
		metrics.Acceleration,
	}
}

func (s *LSM303Accelerometer) Verify() bool {
	if !s.I2C.Verify() {
		return false
	}

	if devID, err := s.I2C.ReadReg(LSM303C_A_DEVICE_ID_REGISTER); err == nil {
		return devID == LSM303C_A_DEVICE_ID
	}

	return false
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
	s.Lock()
	defer s.Unlock()

	x, y, z, err := s.dev.SenseRaw(); if err != nil {
		return model.Vector{}, err
	}

	return model.Vector {
		X: float64(x),
		Y: float64(y),
		Z: float64(z),
	}, nil
}

// ReadTemperature parses data returned as magnetic force vector
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
	ctx.WriterFor(metrics.Magnetism).WriteWithError(toMagnitude(s.ReadAxes()))
	ctx.WriterFor(metrics.Temperature).WriteWithError(s.ReadTemperature())
}

func (s *LSM303Magnetometer) Metrics() []models.Metric {
	return []models.Metric{
		metrics.Magnetism,
		metrics.Temperature,
	}
}

func (s *LSM303Magnetometer) Verify() bool {
	if !s.I2C.Verify() {
		return false
	}

	if devID, err := s.I2C.ReadReg(LSM303C_M_DEVICE_ID_REGISTER); err == nil {
		return devID == LSM303C_M_DEVICE_ID
	}

	return false
}
