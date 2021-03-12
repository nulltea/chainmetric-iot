package device

import (
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/model"
	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

func DiscoverSpecs() (*model.DeviceSpecs, error) {
	var (
		availableMetrics = make(map[models.Metric]bool)
	)

	network, err := shared.GetNetworkEnvironmentInfo(); if err != nil {
		return nil, err
	}

	// for bus, addrs := range shared.ScanI2CAddrs(0x40, 0x76) { // TODO: smart min & max addresses definition
	// 	for _, addr := range addrs {
	// 		if sf, ok := sensors.I2CSensorsMap[addr]; ok {
	// 			for _, metric := range sf(bus).Metrics() {
	// 				availableMetrics[metric] = true
	// 			}
	// 		}
	// 	}
	// }

	shared.Logger.Debug("I2C scan ended")

	var (
		supports = make([]string, len(availableMetrics))
		i = 0
	)

	for m, _ := range availableMetrics {
		supports[i] = string(m)
		i++
	}

	return &model.DeviceSpecs{
		Network: *network,
		Supports: supports,
	}, nil
}
