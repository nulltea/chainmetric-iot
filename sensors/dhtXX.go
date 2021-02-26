package sensors

import (
	"github.com/d2r2/go-dht"
)

type DHTxx struct {
	sensorType dht.SensorType
	pin        int
}

func NewDHTxx(deviceID string, pin int) *DHTxx {
	return &DHTxx{
		sensorType: sensorTypeDHT(deviceID),
		pin:        pin,
	}
}

func (s *DHTxx) Read() (float32, float32, error) {
	temperature, humidity, _, err := dht.ReadDHTxxWithRetry(s.sensorType, s.pin, false, 10)

	return temperature, humidity, err
}

func sensorTypeDHT(deviceID string) dht.SensorType {
	switch deviceID {
	case "DHT11":
		return dht.DHT11
	case "DHT12":
		return dht.DHT12
	case "DHT22":
		return dht.DHT22
	default:
		return dht.DHT11
	}
}
