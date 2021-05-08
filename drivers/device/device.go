package device

import (
	"context"
	"sync"
	"time"

	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/periphery"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
	"github.com/timoth-y/chainmetric-sensorsys/engine"
	"github.com/timoth-y/chainmetric-sensorsys/model"
	"github.com/timoth-y/chainmetric-sensorsys/network/blockchain"
	"github.com/timoth-y/chainmetric-sensorsys/network/local"
)

// Device defines driver for the IoT device itself.
type Device struct {
	ctx context.Context

	specs    *model.DeviceSpecs
	model    *models.Device
	assets   *assetsCache
	requests *requirementsCache

	reader   *engine.SensorsReader
	client   *blockchain.Client
	localnet *local.Client

	detectedI2Cs  periphery.I2CDetectResults
	staticSensors []sensor.Sensor

	pingTimer *time.Timer

	active       bool
	cancelDevice context.CancelFunc
}

// New constructs new IoT Device driver instance.
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

// SetClient assigns the blockchain.Client to the device.
func (d *Device) SetClient(client *blockchain.Client) *Device {
	d.client = client
	return d
}

// SetEngine assigns the engine.SensorsReader engine to the device.
func (d * Device) SetEngine(reader *engine.SensorsReader) *Device {
	d.reader = reader
	return d
}

// SetLocalNet assigns the local.Client for low range communications to the device.
func (d * Device) SetLocalNet(localnet *local.Client) *Device {
	d.localnet = localnet
	return d
}

// RegisterStaticSensors allows to registrant static (not auto-detectable) sensors.
func (d *Device) RegisterStaticSensors(sensors ...sensor.Sensor) *Device {
	d.staticSensors = append(d.staticSensors, sensors...)
	return d
}

func (d *Device) Close() error {
	d.active = false
	d.cancelDevice()

	return nil
}
