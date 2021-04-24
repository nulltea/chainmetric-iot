package device

import (
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/netwrok"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/network"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/periphery"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/sensor"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/sensors"
	"github.com/timoth-y/iot-blockchain-sensorsys/model"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/state"
)

type i2cScanResults map[int][]uint16

func (d *Device) DiscoverSpecs() (*model.DeviceSpecs, error) {
	var (
		availableMetrics = make(map[models.Metric]bool)
	)

	network, err := network.GetNetworkEnvironmentInfo(); if err != nil {
		return nil, err
	}

	d.i2cScan = periphery.DetectI2C(sensors.I2CAddressesRange())

	for bus, addrs := range d.i2cScan {
		for _, addr := range addrs {
			if sf, ok := sensors.LocateI2CSensor(addr); ok {
				for _, metric := range sf.Build(bus).Metrics() {
					availableMetrics[metric] = true
				}
			}
		}
	}

	for _, ss := range d.staticSensors {
		for _, metric := range ss.Metrics() {
			availableMetrics[metric] = true
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
		State: state.Online,
	}

	return d.specs, nil
}

func (d *Device) SupportedSensors() []sensor.Sensor {
	var supports = make([]sensor.Sensor, 0)

	if d.i2cScan == nil {
		d.DiscoverSpecs()
	}

	for bus, addrs := range d.i2cScan {
		for _, addr := range addrs {
			if sf, ok := sensors.LocateI2CSensor(addr); ok {
				supports = append(supports, sf.Build(bus))
			}
		}
	}

	supports = append(supports, d.staticSensors...)

	return supports
}
