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
	ctx        context.Context
	state      *models.Device
	stateMutex sync.Mutex
	specs      model.DeviceSpecs

	cacheLayer

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
		cacheLayer: newCacheLayer(),
		staticSensors: make([]sensor.Sensor, 0),
		cancelDevice: cancel,
	}
}

// ID returns Device unique identifier key in blockchain network.
func (d *Device) ID() string {
	return d.state.ID
}

// UpdateDeviceModel updates Device data model (models.Device).
//
// This method won't change state of the Device in blockchain ledger,
// it just updates device locally saved properties.
func (d *Device) UpdateDeviceModel(model *models.Device) {
	d.stateMutex.Lock()
	defer d.stateMutex.Unlock()

	d.state = model
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

	req := requests.DeviceUpdateRequest{
		State: &state,
	}

	if err := blockchain.Contracts.Devices.Update(d.state.ID, req); err != nil {
		return errors.Wrap(err,"failed to update device state")
	}

	req.Update(d.state)

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

	req := requests.DeviceUpdateRequest{
		Location: &location,
	}

	if err := blockchain.Contracts.Devices.Update(d.state.ID, req); err != nil {
		return errors.Wrap(err,"failed to update device location")
	}

	req.Update(d.state)

	return nil
}

// Specs returns Device current specification.
func (d *Device) Specs() model.DeviceSpecs {
	return d.specs
}

// SetSpecs updates Device current specification in blockchain network.
func (d *Device) SetSpecs(setter func(specs *model.DeviceSpecs)) error {
	specs := &model.DeviceSpecs{}
	*specs = d.specs
	setter(specs)

	if len(specs.Supports) == 0 {
		return errors.New("conflict setting state: device must support at least one metric")
	}

	if len(specs.Hostname) == 0 {
		return errors.New("conflict setting state: hostname must be defines for the device")
	}

	if len(specs.IPAddress) == 0 {
		return errors.New("conflict setting state: IP address must be defines for the device")
	}

	req := requests.DeviceUpdateRequest{
		Hostname: &specs.Hostname,
		IP: &specs.IPAddress,
		Supports: specs.Supports,
	}

	if err := blockchain.Contracts.Devices.Update(d.state.ID, req); err != nil {
		return errors.Wrap(err,"failed to update device specs")
	}

	req.Update(d.state)

	return nil
}

// Battery returns Device current battery stats.
func (d *Device) Battery() models.DeviceBattery {
	return d.state.Battery
}

// SetBattery updates Device current battery stats in blockchain network.
func (d *Device) SetBattery(battery models.DeviceBattery) error {1
	req := requests.DeviceUpdateRequest{
		Battery: &battery,
	}

	if err := blockchain.Contracts.Devices.Update(d.state.ID, req); err != nil {
		return errors.Wrap(err,"failed to update device specs")
	}

	req.Update(d.state)

	return nil
}

// RegisteredSensors returns map with sensors registered on the Device.
func (d *Device) RegisteredSensors() sensor.SensorsRegister {
	return d.reader.RegisteredSensors()
}

// RegisterSensors adds given `sensors` on the Device sensors pool.
func (d *Device) RegisterSensors(sensors ...sensor.Sensor) {
	d.reader.RegisterSensors()
}

// UnregisterSensor removes sensor by given `id` from the Device sensors pool.
func (d *Device) UnregisterSensor(id string) {
	d.reader.UnregisterSensor(id)
}

// SetEngine assigns the engine.SensorsReader engine to the device.
func (d *Device) SetEngine(reader *engine.SensorsReader) *Device {
	d.reader = reader
	return d
}

// StaticSensors returns map with sensors statically registered on the Device.
func (d *Device) StaticSensors() sensor.SensorsRegister {
	sMap := make(map[string]sensor.Sensor)

	for i, sensor := range d.staticSensors {
		sMap[sensor.ID()] = d.staticSensors[i]
	}

	return sMap
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

func (d *Device) NotifyOff() error {
	if d.model == nil {
		return nil
	}

	return blockchain.Contracts.Devices.UpdateState(d.model.ID, state.Offline)
}
