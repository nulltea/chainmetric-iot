package device

import (
	"context"

	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/config"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/display"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/sensors"
	"github.com/timoth-y/iot-blockchain-sensorsys/gateway/blockchain"
	"github.com/timoth-y/iot-blockchain-sensorsys/model"
)

type Device struct {
	Specs *model.DeviceSpecs
	model *models.Device
	assets map[string]bool

	display display.Display
	client  *blockchain.Client
	config  config.Config

	i2cScan       map[int][]uint8
	staticSensors []sensors.Sensor

	cancelEvents context.CancelFunc
}

func NewDevice() *Device {
	return &Device{
		assets: make(map[string]bool),
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

func (d *Device) Close() error {
	if err := d.display.Close(); err != nil {
		return err
	}

	d.client.Close()

	d.cancelEvents()

	return nil
}

