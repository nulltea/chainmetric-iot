package utils

import (
	"encoding/json"

	"github.com/skip2/go-qrcode"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/model"
	"github.com/timoth-y/iot-blockchain-sensorsys/sensors"
)

func GenerateDeviceSignature() (*model.DeviceSignature, error) {
	var (
		availableMetrics = make(map[models.Metric]bool)
	)

	network, err := GetNetworkEnvironmentInfo(); if err != nil {
		return nil, err
	}

	for bus, addrs := range ScanI2CAddrs(0x40, 0x76) {
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

	return &model.DeviceSignature{
		Network: *network,
		Supports: supports,
	}, nil
}

func GenerateDeviceSignatureInQR() ([]byte, error) {
	sig, err := GenerateDeviceSignature(); if err != nil {
		return nil, err
	}

	content, err := json.Marshal(sig); if err != nil {
		return nil, err
	}

	qrcode.WriteFile(string(content), qrcode.Medium, 64, "../qr.png")

	return qrcode.Encode(string(content), qrcode.Medium, 64)
}
