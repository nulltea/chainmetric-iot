package device

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/periphery"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/sensor"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/sensors"
	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

var (
	hotswapOnce = sync.Once{}
)

func (d *Device) initHotswap() {
	hotswapOnce.Do(func() {
		var (
			startTime  time.Time
		)

		ctx, cancel := context.WithCancel(context.Background())
		d.cancelHotswap = func() {
			cancel()
			hotswapOnce = sync.Once{}
		}

		go func() {
		LOOP: for {
			startTime = time.Now()

			if err := d.handleHotswap(); err != nil {
				shared.Logger.Error(errors.Wrap(err, "failed to handle hotswap"))
			}

			select {
			case <-time.After(time.Second - time.Since(startTime)):
			case <- ctx.Done():
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

	d.detectedI2Cs = periphery.DetectI2C(sensors.I2CAddressesRange())

	for bus, addrs := range d.detectedI2Cs {
		for _, addr := range addrs {
			if sf, ok := sensors.LocateI2CSensor(addr); ok {
				s := sf.Build(bus)
				detectedSensors[s.ID()] = s
			}
		}
	}

	for id := range detectedSensors {
		if _, ok := registeredSensors[id]; !ok {
			d.reader.RegisterSensors(detectedSensors[id])
			isChanges = true
		}
	}

	for id := range registeredSensors {
		if _, ok := detectedSensors[id]; !ok {
			d.reader.UnregisterSensor(id)
			isChanges = true
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
