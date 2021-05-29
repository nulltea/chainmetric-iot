package device

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/models/requests"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/device/modules"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
	"github.com/timoth-y/chainmetric-sensorsys/model"
	"github.com/timoth-y/chainmetric-sensorsys/network/blockchain"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// Device defines driver for the IoT device itself.
type Device struct {
	ctx        context.Context
	state      *models.Device
	stateMutex sync.Mutex
	specs      model.DeviceSpecs
	modulesReg modules.Registry

	cacheLayer

	sensors       sensor.SensorsRegister
	staticSensors sensor.SensorsRegister

	active       bool
	cancelDevice context.CancelFunc
}

// New constructs new IoT Device driver instance.
func New(modules ...modules.Module) *Device {
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
	if d.IsLoggedToNetwork() {
		shared.Logger.Warning("won't set state since device hasn't been logged yet")
		return nil
	}

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
	if d.IsLoggedToNetwork() {
		shared.Logger.Warning("won't set location since device hasn't been logged yet")
		return nil
	}

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

	if err := blockchain.Contracts.Devices.Update(d.ID(), req); err != nil {
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

	d.specs = *specs

	req := requests.DeviceUpdateRequest{
		Hostname: &specs.Hostname,
		IP: &specs.IPAddress,
		Supports: specs.Supports,
	}

	if d.IsLoggedToNetwork() {
		if err := blockchain.Contracts.Devices.Update(d.ID(), req); err != nil {
			return errors.Wrap(err,"failed to update device specs")
		}
	} else {
		shared.Logger.Warning("won't update specs since device hasn't been logged yet")
		return nil
	}

	req.Update(d.state)

	return nil
}

// Battery returns Device current battery stats.
func (d *Device) Battery() models.DeviceBattery {
	return d.state.Battery
}

// SetBattery updates Device current battery stats in blockchain network.
func (d *Device) SetBattery(battery models.DeviceBattery) error {
	if d.IsLoggedToNetwork() {
		shared.Logger.Warning("won't set battery info since device hasn't been logged yet")
		return nil
	}

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
	return d.sensors
}

// RegisterSensors adds given `sensors` on the Device sensors pool.
func (d *Device) RegisterSensors(sensors ...sensor.Sensor) {
	for i, s := range sensors {
		d.sensors[s.ID()] = sensors[i]
	}

	if d.IsLoggedToNetwork() {
		d.updateSupportedMetrics()
	}
}

// UnregisterSensor removes sensor by given `id` from the Device sensors pool.
func (d *Device) UnregisterSensor(id string) {
	if d.sensors.Exists(id) {
		delete(d.sensors, id)
	}

	if d.IsLoggedToNetwork() {
		d.updateSupportedMetrics()
	}
}

// UpdateSensorsRegister applies changes in sensor.SensorsRegister of the Device.
func (d *Device) UpdateSensorsRegister(added []sensor.Sensor, removed []string) {
	for i, s := range added {
		d.sensors[s.ID()] = added[i]
	}

	for _, id := range removed {
		delete(d.sensors, id)
	}

	if d.IsLoggedToNetwork() {
		d.updateSupportedMetrics()
	}
}

func (d *Device) updateSupportedMetrics() {
	if err := blockchain.Contracts.Devices.Update(d.ID(), requests.DeviceUpdateRequest{
		Supports: d.sensors.SupportedMetrics(),
	}); err != nil {
		shared.Logger.Error(errors.Wrap(err, "failed to update supported metrics"))
	}
}

// StaticSensors returns map with sensors statically registered on the Device.
func (d *Device) StaticSensors() sensor.SensorsRegister {
	return d.staticSensors
}

// RegisterStaticSensors allows to registrant static (not auto-detectable) sensors.
func (d *Device) RegisterStaticSensors(sensors ...sensor.Sensor) *Device {
	for i, s := range sensors {
		d.staticSensors[s.ID()] = sensors[i]
	}
	return d
}

func (d *Device) Close() error {
	d.active = false
	d.cancelDevice()

	return nil
}
