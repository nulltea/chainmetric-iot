package sensors

import (
	"math"

	"github.com/d2r2/go-i2c"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/model"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
)

// Earth Gravity constant in [m/s^2]
const earthGravityMS2 = 9.80665

// The typical scale factor in g/LSB
const scaleMultiplier = 0.0039

// Registers
const (
	dataFormat = 0x31
	bwRate     = 0x2C
	powerCTL   = 0x2D
	measure    = 0x08
)

// Device bandwidth and output data rates
const (
	Rate1600HZ = 0x0F
	Rate800HZ  = 0x0E
	Rate400HZ  = 0x0D
	Rate200HZ  = 0x0C
	Rate100HZ  = 0x0B
	Rate50HZ   = 0x0A
	Rate25HZ   = 0x09
)

// Measurement Range
const (
	Range2G  = 0x00
	Range4G  = 0x01
	Range8G  = 0x02
	Range16G = 0x03
)

// Axes Data
const (
	dataX0 = 0x32
	dataX1 = 0x33
	dataY0 = 0x34
	dataY1 = 0x35
	dataZ0 = 0x36
	dataZ1 = 0x37
)

// Represents ADXL345 sensor device
type ADXL345 struct {
	addr uint8
	bus int
	i2c *i2c.I2C
}

func NewADXL345(addr uint8, bus int) *ADXL345 {
	return &ADXL345{
		addr: addr,
		bus: bus,
	}
}

func (s *ADXL345) ID() string {
	return "ADXL345"
}

func (s *ADXL345) Init() (err error) {
	s.i2c, err = i2c.NewI2C(s.addr, s.bus); if err != nil {
		return
	}

	// changes the device bandwidth and output data rate
	if err = s.i2c.WriteRegU8(bwRate, Rate100HZ); err != nil {
		return
	}

	if err = s.setRange(Range2G); err != nil {
		return
	}

	// enables measurement on sensor
	if err = s.i2c.WriteRegU8(powerCTL, measure); err != nil {
		return
	}

	return
}

// ReadAxesG retrieves axes acceleration data as multiplications of G
func (s *ADXL345) ReadAxesG() (model.Vector, error) {
	buf, _, err := s.i2c.ReadRegBytes(dataX0, 6); if err != nil {
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
		X: round(gAxes.X*earthGravityMS2, 4),
		Y: round(gAxes.Y*earthGravityMS2, 4),
		Z: round(gAxes.Z*earthGravityMS2, 4),
	}, nil
}

func (s *ADXL345) Harvest(ctx *Context) {
	ctx.For(metrics.AccelerationInG).WriteWithError(s.ReadAxesG())
	ctx.For(metrics.AccelerationInMS2).WriteWithError(s.ReadAxesMS2())
}

func (s *ADXL345) Metrics() []models.Metric {
	return []models.Metric{
		metrics.AccelerationInG,
		metrics.AccelerationInMS2,
	}
}

func (s *ADXL345) Verify() bool {
	return true // TODO verify by device ID
}

func (s *ADXL345) Active() bool {
	return s.i2c != nil
}

// Close disconnects from the device
func (s *ADXL345) Close() error {
	defer s.clean()
	return s.i2c.Close()
}

// setRange changes the range of sensor. Available ranges are 2G, 4G, 8G and 16G.
func (s *ADXL345) setRange(newRange byte) error {
	format, err := s.i2c.ReadRegU8(dataFormat); if err != nil {
		return err
	}

	value := int32(format)
	value &= ^0x0F
	value |= int32(newRange)
	value |= 0x08

	return s.i2c.WriteRegU8(dataFormat, byte(value))
}

func round(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Floor(f*shift+.5) / shift
}

func (s *ADXL345) clean() {
	s.i2c = nil
}
