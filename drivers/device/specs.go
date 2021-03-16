package device

import (
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/sensors"
	"github.com/timoth-y/iot-blockchain-sensorsys/model"
	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

type i2cScanResults map[int][]uint8

func (d *Device) DiscoverSpecs() (*model.DeviceSpecs, error) {
	var (
		availableMetrics = make(map[models.Metric]bool)
	)

	network, err := shared.GetNetworkEnvironmentInfo(); if err != nil {
		return nil, err
	}

	d.i2cScan = shared.ScanI2CAddrs(0x40, 0x76) // TODO: smart min & max addresses definition

	for bus, addrs := range d.i2cScan {
		for _, addr := range addrs {
			if sf, ok := sensors.I2CSensorsMap[addr]; ok {
				for _, metric := range sf(bus).Metrics() {
					availableMetrics[metric] = true
				}
			}
		}
	}

	var (
		supports = make([]string, len(availableMetrics))
		i = 0
	)

	for m, _ := range availableMetrics {
		supports[i] = string(m)
		i++
	}

	d.specs = &model.DeviceSpecs{
		Network: *network,
		Supports: supports,
	}

	return d.specs, nil
}

func (d *Device) SupportedSensors() []sensors.Sensor {
	var supports = make([]sensors.Sensor, 0)

	if d.i2cScan == nil {
		d.DiscoverSpecs()
	}

	for bus, addrs := range d.i2cScan {
		for _, addr := range addrs {
			if sf, ok := sensors.I2CSensorsMap[addr]; ok {
				supports = append(supports, sf(bus))
			}
		}
	}

	supports = append(supports, d.staticSensors...)

	return supports
}
