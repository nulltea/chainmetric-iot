package periphery

import "sync"

// An I2COption configures a ADC driver.
type I2COption interface {
	Apply(adc *I2C)
}

// I2COptionFunc is a function that configures a I2C driver.
type I2COptionFunc func(d *I2C)

// Apply calls I2COptionFunc on the driver instance.
func (f I2COptionFunc) Apply(i2c *I2C) {
	f(i2c)
}

// WithMutex can be used to setup mutex I2C driver.
// Default is a new sync.Mutex instance.
func WithMutex(mutex *sync.Mutex) I2COption {
	return I2COptionFunc(func(d *I2C) {
		d.Mutex = mutex
	})
}
