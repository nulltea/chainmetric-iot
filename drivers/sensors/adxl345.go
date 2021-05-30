package sensors

import (
	"math"
	"sync"

	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-core/models/metrics"

	"github.com/timoth-y/chainmetric-sensorsys/core/dev/sensor"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/periphery"
	"github.com/timoth-y/chainmetric-sensorsys/model"
)

var (
	adxl345Mutex = &sync.Mutex{}
)

const (
	// The typical scale factor in g/LSB
	scaleMultiplier = 0.0039
)

// ADXL345 sensor device.
type ADXL345 struct {
	*periphery.I2C
}

func NewADXL345(addr uint16, bus int) sensor.Sensor {
	return &ADXL345{
		I2C: periphery.NewI2C(addr, bus, periphery.WithMutex(adxl345Mutex)),
	}
}

func (s *ADXL345) ID() string {
	return "ADXL345"
}

func (s *ADXL345) Init() error {
	if err := s.I2C.Init(); err != nil {
		return err
	}

	// changes the device bandwidth and output data rate
	if err := s.WriteRegBytes(ADXL345_BW_RATE, ADXL345_Rate100HZ); err != nil {
		return err
	}

	if err := s.setRange(ADXL345_RANGE2G); err != nil {
		return err
	}

	// enables measurement on sensor
	if err := s.WriteRegBytes(ADXL345_POWER_CTL, ADXL345_MEASURE); err != nil {
		return err
	}

	return nil
}

// ReadAxes retrieves axes acceleration data as multiplications of G
func (s *ADXL345) ReadAxes() (model.Vector, error) {
	buf, err := s.ReadRegBytes(ADXL345_DATAX0, 6); if err != nil {
		return model.Vector{}, err
	}

	x := int16(buf[0]) | (int16(buf[1]) << 8)
	y := int16(buf[2]) | (int16(buf[3]) << 8)
	z := int16(buf[4]) | (int16(buf[5]) << 8)

	return model.Vector {
		X: round(float64(x) * scaleMultiplier, 4),
		Y: round(float64(y) * scaleMultiplier, 4),
		Z: round(float64(z) * scaleMultiplier, 4),
	}, nil
}

func (s *ADXL345) Harvest(ctx *sensor.Context) {
	ctx.WriterFor(metrics.Acceleration).WriteWithError(toMagnitude(s.ReadAxes()))
}

func (s *ADXL345) Metrics() []models.Metric {
	return []models.Metric{
		metrics.Acceleration,
	}
}

func (s *ADXL345) Verify() bool {
	if !s.I2C.Verify() {
		return false
	}

	if devID, err := s.I2C.ReadReg(ADXL345_DEVICE_ID_REGISTER); err == nil {
		return devID == ADXL345_DEVICE_ID
	}

	return false
}

// setRange changes the range of sensor. Available ranges are 2G, 4G, 8G and 16G.
func (s *ADXL345) setRange(newRange byte) error {
	format, err := s.ReadReg(ADXL345_DATA_FORMAT); if err != nil {
		return err
	}

	value := int32(format)
	value &= ^0x0F
	value |= int32(newRange)
	value |= 0x08

	return s.WriteRegBytes(ADXL345_DATA_FORMAT, byte(value))
}

func round(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Floor(f*shift+.5) / shift
}

func toMagnitude(vector model.Vector, err error) (float64, error) {
	r := math.Pow(vector.X, 2) + math.Pow(vector.Y, 2) + math.Pow(vector.Z, 2)

	return math.Sqrt(r), err
}
