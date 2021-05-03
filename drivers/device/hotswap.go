package device

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/periphery"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensors"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

var (
	hotswapOnce = sync.Once{}
)

func (d *Device) initHotswap() {
	if !d.active {
		return
	}

	hotswapOnce.Do(func() {
		var (
			startTime  time.Time
			interval = viper.GetDuration("device.hotswap_detect_interval")
		)

		go func() {
			LOOP: for {
				startTime = time.Now()

				if err := d.handleHotswap(); err != nil {
					shared.Logger.Error(errors.Wrap(err, "failed to handle hotswap"))
				}

				select {
				case <-time.After(interval - time.Since(startTime)):
				case <- d.ctx.Done():
					shared.Logger.Debug("Hotswap routine ended.")
					break LOOP
				}
			}
		}()
	})
}

func (d *Device) handleHotswap() error {
	var (
		detectedSensors = make(map[string]sensor.Sensor)
		registeredSensors = d.reader.RegisteredSensors()
		contract = d.client.Contracts.Devices
		isChanges bool
	)

	d.detectedI2Cs = periphery.ScanI2C(sensors.I2CAddressesRange(), sensors.LocateI2CSensor)
	for _, devices := range d.detectedI2Cs {
		for _, s := range devices {
			detectedSensors[s.ID()] = s
		}
	}

	for id := range registeredSensors {
		if _, ok := detectedSensors[id]; !ok && !d.isStaticSensor(id) {
			d.reader.UnregisterSensor(id)
			isChanges = true
			shared.Logger.Debugf("Hotswap: %s sensor was detached from the device", id)
		}
	}

	for id := range detectedSensors {
		if _, ok := registeredSensors[id]; !ok {
			d.reader.RegisterSensors(detectedSensors[id])
			isChanges = true
			shared.Logger.Debugf("Hotswap: %s sensor was attached to the device", id)
		}
	}

	if isChanges {
		if _, err := d.DiscoverSpecs(false); err != nil {
			return err
		}

		return contract.UpdateSpecs(d.model.ID, d.specs)
	}

	return nil
}

func (d *Device) isStaticSensor(id string) bool {
	for i := range d.staticSensors {
		if d.staticSensors[i].ID() == id {
			return true
		}
	}

	return false
}
