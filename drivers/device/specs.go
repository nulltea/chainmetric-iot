package device

import (
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/network"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/periphery"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensors"
	"github.com/timoth-y/chainmetric-sensorsys/model"
	"github.com/timoth-y/chainmetric-sensorsys/model/state"
)

func (d *Device) DiscoverSpecs(rescan bool) (*model.DeviceSpecs, error) {
	var (
		availableMetrics = make(map[models.Metric]bool)
	)

	network, err := network.GetNetworkEnvironmentInfo(); if err != nil {
		return nil, err
	}

	if rescan || d.detectedI2Cs == nil {
		d.detectedI2Cs = periphery.DetectI2C(sensors.I2CAddressesRange())
	}

	for bus, addrs := range d.detectedI2Cs {
		for _, addr := range addrs {
			if sf, ok := sensors.LocateI2CSensor(addr, bus); ok {
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

	d.DiscoverSpecs(d.detectedI2Cs == nil)

	for bus, addrs := range d.detectedI2Cs {
		for _, addr := range addrs {
			if sf, ok := sensors.LocateI2CSensor(addr, bus); ok {
				supports = append(supports, sf.Build(bus))
			}
		}
	}

	supports = append(supports, d.staticSensors...)

	return supports
}
