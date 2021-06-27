package device

import (
	"context"
	"sync"

	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-iot/core/dev/sensor"
	"github.com/timoth-y/chainmetric-iot/model"
)

// Device defines driver for the IoT device itself.
type Device struct {
	ctx        context.Context
	state      *models.Device
	stateMutex sync.Mutex
	specs      model.DeviceSpecs
	modulesReg ModulesRegistry

	cacheLayer

	sensors       sensor.SensorsRegister
	staticSensors sensor.SensorsRegister

	active       bool
	cancelDevice context.CancelFunc
}

// New constructs new IoT Device driver instance.
func New(modules ...Module) *Device {
	ctx, cancel := context.WithCancel(context.Background())

	dev := &Device{
		ctx:           ctx,
		cacheLayer:    newCacheLayer(),
		sensors:       make(sensor.SensorsRegister),
		staticSensors: make(sensor.SensorsRegister),
		cancelDevice:  cancel,
	}

	dev.modulesReg = modules

	return dev
}

// ID returns Device unique identifier key in blockchain network.
func (d *Device) ID() string {
	if !d.IsLoggedToNetwork() {
		return ""
	}

	return d.state.ID
}

// Name returns Device given name.
func (d *Device) Name() string {
	if !d.IsLoggedToNetwork() {
		return ""
	}

	return d.state.Name
}

// IsLoggedToNetwork determines whether the Device is logged to network and thus is ready to operate.
func (d *Device) IsLoggedToNetwork() bool {
	return d.state != nil
}

// Start performs startup of the Device and setting up all registered modules.
func (d *Device) Start() {
	d.modulesReg.Setup(d)
	d.modulesReg.Start(d.ctx)
}

// Close stops all working device.Module and frees allocated resources.
func (d *Device) Close() error {
	d.active = false
	d.cancelDevice()
	d.modulesReg.Close()

	return nil
}
