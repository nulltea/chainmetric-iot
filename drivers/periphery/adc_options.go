package periphery

import "sync"

// An ADCOption configures a ADC driver.
type ADCOption interface {
	Apply(adc *ADS1115)
}

// ADCOptionFunc is a function that configures a ADC driver.
type ADCOptionFunc func(d *ADS1115)

// Apply calls ADCOptionFunc on the driver instance.
func (f ADCOptionFunc) Apply(adc *ADS1115) {
	f(adc)
}

// WithConversion can be used to setup ADC readings conversion.
// Default is a function that returns input value as is.
func WithConversion(convertor func(v float64) float64) ADCOption {
	return ADCOptionFunc(func(d *ADS1115) {
		d.convertor = convertor
	})
}

// WithBias can be used to specify ADC readings bias.
// Default is 0.
func WithBias(bias float64) ADCOption {
	return ADCOptionFunc(func(d *ADS1115) {
		d.bias = bias
	})
}

// WithI2CMutex can be used to specify mutex for I2C bus driver.
// Default is a new sync.Mutex instance.
func WithI2CMutex(mutex *sync.Mutex) ADCOption {
	return ADCOptionFunc(func(d *ADS1115) {
		d.Mutex = mutex
	})
}
