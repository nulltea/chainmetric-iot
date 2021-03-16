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

	i2cScan i2cScan
	staticSensors []sensors.Sensor

	cancelEvents context.CancelFunc
}

func NewDevice() *Device {
	return &Device{
		assets: &assetsCache{
			mutex: sync.Mutex{},
			data:  make(map[string]bool),
		},
		requests: &requirementsCache{
			mutex: sync.Mutex{},
			data:  make(map[string]*readingsRequest),
		},
		staticSensors: make([]sensors.Sensor, 0),
	}
}

func (d *Device) SetConfig(cfn config.Config) *Device {
	d.config = cfn
	return d
}

func (d *Device) RegisterStaticSensors(sensors ...sensors.Sensor) *Device {
	d.staticSensors = append(d.staticSensors, sensors...)
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

	d.cancelEvents()

	d.reader.Close()

	d.client.Close()

	return nil
}
