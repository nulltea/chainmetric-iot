package sensor

import "github.com/timoth-y/chainmetric-sensorsys/core"

// Factory defines interface for building core.Sensor.
type Factory interface {
	Build(bus int) core.Sensor
}

// FactoryFunc builds core.Sensor.
type FactoryFunc func(int) core.Sensor

// Build calls FactoryFunc to build core.Sensor on specified peripheral bus.
func (f FactoryFunc) Build(bus int) core.Sensor {
	return f(bus)
}

// I2CFactory provides new factory for building I2C-based core.Sensor.
func I2CFactory(factory func(addr uint16, bus int) core.Sensor, addr uint16) Factory {
	return FactoryFunc(func(bus int) core.Sensor {
		return factory(addr, bus)
	})
}
