package sensors

import (
	"github.com/d2r2/go-dht"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/engine/sensor"
	"github.com/timoth-y/iot-blockchain-sensorsys/model/metrics"
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

func NewDHT11(pin int) *DHTxx {
	return &DHTxx{
		sensorType: dht.DHT11,
		pin:        pin,
	}
}

func NewDHT22(pin int) *DHTxx {
	return &DHTxx{
		sensorType: dht.DHT22,
		pin:        pin,
	}
}

func (s *DHTxx) ID() string {
	return s.sensorType.String()
}

func (s *DHTxx) Init() error {
	return nil
}

func (s *DHTxx) Harvest(ctx *sensor.Context) {
	temperature, humidity, _, err := dht.ReadDHTxxWithRetry(s.sensorType, s.pin, false, 10)

	ctx.For(metrics.Temperature).Write(temperature)
	ctx.For(metrics.Humidity).Write(humidity)
	ctx.Error(err)
}

func (s *DHTxx) Metrics() []models.Metric {
	return []models.Metric {
		metrics.Temperature,
		metrics.Humidity,
	}
}

func (s *DHTxx) Active() bool {
	return true
}

func (s *DHTxx) Close() error {
	return nil
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