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
)

// Device defines driver for the IoT device itself.
type Device struct {
	ctx context.Context

	specs    *model.DeviceSpecs
	model    *models.Device
	assets   *assetsCache
	requests *requirementsCache

	reader   *engine.SensorsReader

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


// SetEngine assigns the engine.SensorsReader engine to the device.
func (d * Device) SetEngine(reader *engine.SensorsReader) *Device {
	d.reader = reader
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
