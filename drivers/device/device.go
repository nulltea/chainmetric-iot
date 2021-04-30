package device

import (
	"context"
	"sync"
	"time"

	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/display"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/periphery"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
	"github.com/timoth-y/chainmetric-sensorsys/engine"
	"github.com/timoth-y/chainmetric-sensorsys/gateway/blockchain"
	"github.com/timoth-y/chainmetric-sensorsys/model"
)

type Device struct {
	ctx context.Context

	specs    *model.DeviceSpecs
	model    *models.Device
	assets   *assetsCache
	requests *requirementsCache

	reader  *engine.SensorsReader
	client  *blockchain.Client
	display display.Display

	detectedI2Cs  periphery.I2CDetectResults
	staticSensors []sensor.Sensor

	pingTimer *time.Timer

	active bool
	cancelDevice context.CancelFunc
}

func New() *Device {
	ctx, cancel := context.WithCancel(context.Background())

	return &Device{
		ctx: ctx,
		assets: &assetsCache{
			mutex: sync.Mutex{},
			data:  make(map[string]bool),
		},
		requests: &requirementsCache{
			mutex: sync.Mutex{},
			data:  make(map[string]*readingsRequest),
		},
		staticSensors: make([]sensor.Sensor, 0),
		cancelDevice: cancel,
	}
}

func (d *Device) RegisterStaticSensors(sensors ...sensor.Sensor) *Device {
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

func (d * Device) DisplayAvailable() bool {
	return d.display != nil && d.display.Active()
}

func (d *Device) Close() error {
	d.active = false
	d.cancelDevice()

	return nil
}
