package periphery

import (
	"math"
	"sort"
	"time"

	"github.com/MichaelS11/go-ads"
	"github.com/pkg/errors"

	"github.com/timoth-y/chainmetric-iot/shared"
)

// ADS1115 ADC chip constants.
const (
	// Conversion constants
	ADS1115_SAMPLES_PER_READ = 32767.0
	ADS1115_VOLTS_PER_SAMPLE = 5

	ADS1115_DEVICE_ID_REGISTER = 0x0E
	ADS1115_DEVICE_ID          = 0x80
)

// ADC defines analog to digital peripheral interface.
type ADC interface {
	// Init performs ADC driver initialisation.
	Init() error
	// Read returns single analog sensor reading value.
	Read() float64
	// RMS aggregates `n` readings from analog sensor and calculates Root Mean Square function.
	RMS(n int, t *time.Duration) float64
	// Max returns max value from `n` analog sensor readings.
	Max(n int, t *time.Duration) float64
	// Min returns min value from `n` analog sensor readings.
	Min(n int, t *time.Duration) float64
	// Verify identifies ADC device and checks it according to implemented driver.
	Verify() bool
	// Active determines whether the ADC device is active.
	Active() bool
	// Close closes connection to ADC device.
	Close() error
}

// ADS1115 implements ADC driver for ADS1115 device.
type ADS1115 struct {
	*ads.ADS
	*I2C
	Addr   uint16
	Bus    string
	active bool

	bias float64
	convertor func(float64) float64
}

// NewADC constructs a new ADC implementation via ADS1115 device driver.
func NewADC(addr uint16, bus int, options ...ADCOption) *ADS1115 {
	d := &ADS1115{
		Bus: shared.NtoI2cBusName(bus),
		Addr: addr,
		I2C: NewI2C(addr, bus),

		convertor: func(v float64) float64 {
			return v
		},
	}

	for i := range options {
		options[i].Apply(d)
	}

	return d
}

// Init sets up the device for communication.
func (d *ADS1115) Init() (err error) {
	if d.ADS, err = ads.NewADS(d.Bus, d.Addr, "ADS1115"); err != nil {
		return errors.Wrapf(err, "failed to init ADS1115 device on '%s' bus and 0x%X address", d.Bus, d.Addr)
	}

	d.active = true

	return nil
}

func (d *ADS1115) Read() float64 {
	d.Lock()
	defer d.Unlock()

	if v, err := d.ADS.ReadRetry(5); err != nil {
		return 0
	} else {
		return d.convertor(float64(v)) - d.bias
	}
}

func (d *ADS1115) RMS(n int, t *time.Duration) float64 {
	var (
		sum float64
		i = n
	)

	for i > 0 {
		if v, err := d.ADS.Read(); err != nil {
			continue
		} else {
			sum +=  math.Pow(float64(v), 2)
		}

		if t != nil {
			time.Sleep(*t)
		}

		i--
	}

	return d.convertor(math.Sqrt(sum / float64(n))) - d.bias
}

func (d *ADS1115) Max(n int, t *time.Duration) float64 {
	results := d.rawSequence(n, t)

	sort.Ints(results)

	return d.convertor(float64(results[len(results) - 1])) - d.bias
}

func (d *ADS1115) Min(n int, t *time.Duration) float64 {
	results := d.rawSequence(n, t)

	sort.Ints(results)

	return d.convertor(float64(results[0])) - d.bias
}

func (d ADS1115) rawSequence(n int, t *time.Duration) []int {
	var (
		i = n
		results []int
	)

	for i > 0 {
		if v, err := d.ADS.Read(); err != nil {
			continue
		} else {
			results = append(results, int(v))
		}

		if t != nil {
			time.Sleep(*t)
		}

		i--
	}

	return results
}

func (d *ADS1115) Verify() bool {
	if !d.I2C.Verify() {
		return false
	}

	if devID, err := d.ReadReg(ADS1115_DEVICE_ID_REGISTER); err == nil {
		return devID == ADS1115_DEVICE_ID
	}

	return false
}

func (d *ADS1115) Active() bool {
	return d.active
}

func (d *ADS1115) Close() error {
	d.active = false
	return d.ADS.Close()
}
