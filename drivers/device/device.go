package device

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/models/requests"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/periphery"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
	"github.com/timoth-y/chainmetric-sensorsys/engine"
	"github.com/timoth-y/chainmetric-sensorsys/model"
	"github.com/timoth-y/chainmetric-sensorsys/network/blockchain"
)

// Device defines driver for the IoT device itself.
type Device struct {
	ctx context.Context

	state    *models.Device
	specs    model.DeviceSpecs
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

// ID returns Device unique identifier key in blockchain network.
func (d *Device) ID() string {
	return d.state.ID
}

// State returns Device current state.
func (d *Device) State() models.DeviceState {
	return d.state.State
}

// SetState updates Device current state in blockchain network.
func (d *Device) SetState(state models.DeviceState) error {
	if d.state.State == state {
		return errors.Errorf("conflict setting state: device state is already '%s'", d.state.State)
	}

	if err := blockchain.Contracts.Devices.Update(d.state.ID, requests.DeviceUpdateRequest{
		State: &state,
	}); err != nil {
		return errors.Wrap(err,"failed to update device state")
	}

	return nil
}

// Location returns Device current location.
func (d *Device) Location() models.Location {
	return d.state.Location
}

// SetLocation updates Device current location in blockchain network.
func (d *Device) SetLocation(location models.Location) error {
	if d.state.Location == location {
		return errors.Errorf(
			"conflict setting location: device location is already '%s'",
			d.state.Location.Name,
		)
	}

	if location.Latitude == 0 || location.Longitude == 0 {
		return errors.New("conflict setting state: new location must contains both coordinates")
	}

	if err := blockchain.Contracts.Devices.Update(d.state.ID, requests.DeviceUpdateRequest{
		Location: &location,
	}); err != nil {
		return errors.Wrap(err,"failed to update device location")
	}

	return nil
}

// Specs returns Device current specification.
func (d *Device) Specs() model.DeviceSpecs {
	return d.specs
}

// SetSpecs updates Device current specification in blockchain network.
func (d *Device) SetSpecs(specs model.DeviceSpecs) error {
	if len(specs.Supports) == 0 {
		return errors.New("conflict setting state: device must support at least one metric")
	}

	if len(specs.Hostname) == 0 {
		return errors.New("conflict setting state: hostname must be defines for the device")
	}

	if len(specs.IPAddress) == 0 {
		return errors.New("conflict setting state: IP address must be defines for the device")
	}

	if err := blockchain.Contracts.Devices.Update(d.state.ID, requests.DeviceUpdateRequest{
		Hostname: &specs.Hostname,
		IP: &specs.IPAddress,
		Supports: specs.Supports,
	}); err != nil {
		return errors.Wrap(err,"failed to update device specs")
	}

	return nil
}

// Battery returns Device current battery stats.
func (d *Device) Battery() models.DeviceBattery {
	return d.state.Battery
}

// SetBattery updates Device current battery stats in blockchain network.
func (d *Device) SetBattery(battery models.DeviceBattery) error {
	if err := blockchain.Contracts.Devices.Update(d.state.ID, requests.DeviceUpdateRequest{
		Battery: &battery,
	}); err != nil {
		return errors.Wrap(err,"failed to update device specs")
	}

	return nil
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
