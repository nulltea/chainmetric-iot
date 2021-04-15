package sensors

import (
	"math"

	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/peripherals"
	"github.com/timoth-y/iot-blockchain-sensorsys/model"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
)

const (
	// Earth Gravity constant in [m/s^2]
	earthGravityMS2 = 9.80665

	// The typical scale factor in g/LSB
	scaleMultiplier = 0.0039
)

// Represents ADXL345 sensor device
type ADXL345 struct {
	addr uint8
	bus  int
	dev  *peripherals.I2C
}

func NewADXL345(addr uint16, bus int) *ADXL345 {
	return &ADXL345{
		dev: peripherals.NewI2C(addr, bus),
	}
}

func (s *ADXL345) ID() string {
	return "ADXL345"
}

func (s *ADXL345) Init() (err error) {
	if err = s.dev.Init(); err != nil {
		return
	}

	// changes the device bandwidth and output data rate
	if err = s.dev.WriteRegBytes(ADXL345_BW_RATE, ADXL345_Rate100HZ); err != nil {
		return
	}

	if err = s.setRange(ADXL345_RANGE2G); err != nil {
		return
	}

	// enables measurement on sensor
	if err = s.dev.WriteRegBytes(ADXL345_POWER_CTL, ADXL345_MEASURE); err != nil {
		return
	}

	return
}

// ReadAxesG retrieves axes acceleration data as multiplications of G
func (s *ADXL345) ReadAxesG() (model.Vector, error) {
	buf, err := s.dev.ReadRegBytes(ADXL345_DATAX0, 6); if err != nil {
		return model.Vector{}, err
	}

	x := int16(buf[0]) | (int16(buf[1]) << 8)
	y := int16(buf[2]) | (int16(buf[3]) << 8)
	z := int16(buf[4]) | (int16(buf[5]) << 8)

	return model.Vector {
		X: round(float64(x) *scaleMultiplier, 4),
		Y: round(float64(y) *scaleMultiplier, 4),
		Z: round(float64(z) *scaleMultiplier, 4),
	}, nil
}

// ReadAxesMS2 parses data returned by GetAxesG and returns them in [m/s^2]
func (s *ADXL345) ReadAxesMS2() (model.Vector, error) {
	gAxes, err := s.ReadAxesG(); if err != nil {
		return model.Vector{}, err
	}

	return model.Vector {
		X: round(gAxes.X * earthGravityMS2, 4),
		Y: round(gAxes.Y * earthGravityMS2, 4),
		Z: round(gAxes.Z * earthGravityMS2, 4),
	}, nil
}

func (s *ADXL345) Harvest(ctx *Context) {
	ctx.For(metrics.AccelerationInG).WriteWithError(toMagnitude(s.ReadAxesG()))
	ctx.For(metrics.AccelerationInMS2).WriteWithError(toMagnitude(s.ReadAxesMS2()))
}

func (s *ADXL345) Metrics() []models.Metric {
	return []models.Metric{
		metrics.AccelerationInG,
		metrics.AccelerationInMS2,
	}
}

func (s *ADXL345) Verify() bool {
	return true
}

func (s *ADXL345) Active() bool {
	return s.dev.Active()
}

// Close disconnects from the device
func (s *ADXL345) Close() error {
	return s.dev.Close()
}

// setRange changes the range of sensor. Available ranges are 2G, 4G, 8G and 16G.
func (s *ADXL345) setRange(newRange byte) error {
	format, err := s.dev.ReadReg(ADXL345_DATA_FORMAT); if err != nil {
		return err
	}

	value := int32(format)
	value &= ^0x0F
	value |= int32(newRange)
	value |= 0x08

	return s.dev.WriteRegBytes(ADXL345_DATA_FORMAT, byte(value))
}

func round(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Floor(f*shift+.5) / shift
}

func toMagnitude(vector model.Vector, err error) (float64, error) {
	r := math.Pow(vector.X, 2) + math.Pow(vector.Y, 2) + math.Pow(vector.Z, 2)

	return math.Sqrt(r), err
}
