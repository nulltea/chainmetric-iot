package network

import (
	"net"
	"os"
	"strings"

	"github.com/timoth-y/iot-blockchain-sensorsys/model"
)

func GetNetworkEnvironmentInfo() (*model.Network, error) {
	var (
		ipAddress    string
		macAddress   string
		hostname     string
		hardwareName string
	)

	addresses, err := net.InterfaceAddrs(); if err != nil {
		return nil, err
	}

	for _, address := range addresses {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipAddress = ipnet.IP.String()
			}
		}
	}

	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		if addrs, err := inter.Addrs(); err == nil {
			for _, addr := range addrs {
				if strings.Contains(addr.String(), ipAddress) {
					hardwareName = inter.Name
				}
			}
		}
	}

	netInterface, err := net.InterfaceByName(hardwareName); if err != nil {
		return nil, err
	}

	macAddress = netInterface.HardwareAddr.String()

	if hostname, err = os.Hostname(); err != nil {
		return nil, err
	}

	return &model.Network{
		IPAddress: ipAddress,
		MACAddress: macAddress,
		Hostname: hostname,
	}, nil
}
