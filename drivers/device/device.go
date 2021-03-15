package device

import (
	"context"
	"sync"

	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/config"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/display"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/sensors"
	"github.com/timoth-y/iot-blockchain-sensorsys/engine"
	"github.com/timoth-y/iot-blockchain-sensorsys/gateway/blockchain"
	"github.com/timoth-y/iot-blockchain-sensorsys/model"
)

type Device struct {
	specs    *model.DeviceSpecs
	model    *models.Device
	assets   *assetsCache
	requests *requirementsCache

	reader *engine.SensorsReader

	client  *blockchain.Client
	display display.Display
	config  config.Config

	i2cScan       map[int][]uint8
	staticSensors []sensors.Sensor

	cancelEvents context.CancelFunc
}

type assetsCache struct {
	mutex sync.Mutex
	data  map[string]bool
}

type requirementsCache struct {
	mutex sync.Mutex
	data  map[string]models.Metrics
}

func NewDevice() *Device {
	return &Device{
		assets: &assetsCache{
			mutex: sync.Mutex{},
			data:  make(map[string]bool),
		},
		requests: &requirementsCache{
			mutex: sync.Mutex{},
			data:  make(map[string]models.Metrics),
		},
		staticSensors: make([]sensors.Sensor, 0),
	}
}

func (d *Device) SetConfig(cfn config.Config) *Device {
	d.config = cfn
	return d
}

func (d *Device) SetDisplay(dp display.Display) *Device {
	d.display = dp
	return d
}

func (d *Device) SetClient(client *blockchain.Client) *Device {
	d.client = client
	return d
}

func (d * Device) SetReader(reader *engine.SensorsReader) *Device {
	d.reader = reader
	return d
}

func (d *Device) Close() error {
	if err := d.display.Close(); err != nil {
		return err
	}

	d.client.Close()

	d.cancelEvents()

	return nil
}

func (ac *assetsCache) Get() []string {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()

	var (
		ids = make([]string, len(ac.data))
		i = 0
	)

	for id := range ac.data {
		ids[i] = id
		i++
	}
	return ids
}
