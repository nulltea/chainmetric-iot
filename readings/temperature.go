package readings

import (
	"context"
	"fmt"

	"github.com/d2r2/go-dht"

	"sensorsys/model"
)

func (s *SensorsReader) SubscribeToTemperatureReadings(sensor string, pin int) error {
	s.subscribe(func(ctx context.Context) {
		temperature, humidity := readDHT(sensorTypeDHT(sensor), pin)
		s.readings <- model.MetricReadings{
			model.Temperature: temperature,
			model.Humidity: humidity,
		}
	})

	return nil
}

func readDHT(sensor dht.SensorType, pin int) (float32, float32) {
	temperature, humidity, _, err := dht.ReadDHTxxWithRetry(sensor, pin, false, 10)
	if err != nil {
		fmt.Println(err)
	}
	return temperature, humidity
}

func sensorTypeDHT(sensor string) dht.SensorType {
	switch sensor {
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