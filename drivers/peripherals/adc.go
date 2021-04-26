package peripherals

import (
	"math"
	"sort"
	"time"

	"github.com/MichaelS11/go-ads"
	"github.com/pkg/errors"

	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// ADC defines analog to digital peripheral interface.
type ADC interface {
	Init() error
	Read() (uint16, error)
	ReadRetry(int) (uint16, error)
	Aggregate(n int, t *time.Duration) float64
	Max(n int, t *time.Duration) uint16
	Min(n int, t *time.Duration) uint16
	Sequence(n int, t *time.Duration) []int
	Active() bool
	Close() error
}

// ADS1115 implements ADC driver for ADS1115 device.
type ADS1115 struct {
	*ads.ADS
	Addr   uint16
	Bus    string
	active bool
	bias uint16
}

// NewADC returns a new ADC implementation via ADS1115 device driver.
func NewADC(addr uint16, bus int) *ADS1115 {
	return &ADS1115{
		Bus: shared.NtoI2cBusName(bus),
		Addr: addr,
	}
}

// Init sets up the device for communication.
func (d *ADS1115) Init() (err error) {
	if d.ADS, err = ads.NewADS(d.Bus, d.Addr, "ADS1115"); err != nil {
		return errors.Wrapf(err, "failed to init ADS1115 device on '%s' bus and 0x%X address", d.Bus, d.Addr)
	}

	d.active = true

	return nil
}

func (d *ADS1115) Aggregate(n int, t *time.Duration) float64 {
	var (
		sum float64
		i = n
	)

	for i > 0 {
		if v, err := d.Read(); err != nil {
			continue
		} else {
			sum +=  math.Pow(float64(v), 2)
		}

		if t != nil {
			time.Sleep(*t)
		}

		i--
	}

	return math.Sqrt(sum / float64(n))
}

func (d *ADS1115) Max(n int, t *time.Duration) uint16 {
	results := d.Sequence(n, t)

	sort.Ints(results)

	return uint16(results[len(results) - 1])
}

func (d *ADS1115) Min(n int, t *time.Duration) uint16 {
	results := d.Sequence(n, t)

	sort.Ints(results)

	return uint16(results[0])
}

func (d ADS1115) Sequence(n int, t *time.Duration) []int {
	var (
		i = n
		results []int
	)

	for i > 0 {
		if v, err := d.Read(); err != nil {
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

func (d *ADS1115) Active() bool {
	return d.active
}

func (d *ADS1115) Close() error {
	d.active = false
	return d.ADS.Close()
}
