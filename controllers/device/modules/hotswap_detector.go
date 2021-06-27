package modules

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-iot/controllers/device"
	"github.com/timoth-y/chainmetric-iot/core/dev/sensor"
	"github.com/timoth-y/chainmetric-iot/core/io"
	"github.com/timoth-y/chainmetric-iot/drivers/sensors"
	"github.com/timoth-y/chainmetric-iot/model/events"
	"github.com/timoth-y/chainmetric-iot/shared"
	"github.com/timoth-y/go-eventdriver"
)

// HotswapDetector implements device.Module ofr detecting changes in connected to device.Device sensors.
type HotswapDetector struct {
	moduleBase

	detectedI2Cs io.I2CDetectResults
}

// WithHotswapDetector can be used to setup HotswapDetector logical device.Module onto the device.Device.
func WithHotswapDetector() device.Module {
	return &HotswapDetector{
		moduleBase: withModuleBase("HOTSWAP_DETECTOR"),
	}
}

func (m *HotswapDetector) Start(ctx context.Context) {
	go m.Do(func() {
		var (
			interval = viper.GetDuration("device.hotswap_detect_interval")
			startTime  time.Time
		)

	LOOP:
		for {
			startTime = time.Now()

			if err := m.handleHotswap(ctx); err != nil {
				shared.Logger.Error(errors.Wrap(err, "failed to handle hotswap"))
			}

			select {
			case <- time.After(interval - time.Since(startTime)):
			case <- ctx.Done():
				shared.Logger.Debug("Hotswap detector module routine ended")
				break LOOP
			}
		}
	})
}

func (m *HotswapDetector) handleHotswap(ctx context.Context) error {
	var (
		detectedSensors = make(sensor.SensorsRegister)
		registeredSensors = m.RegisteredSensors()
		staticSensors = m.StaticSensors()
		payload = events.SensorsRegisterChangedPayload{}
		isChanges bool
	)

	m.detectedI2Cs = io.ScanI2C(sensors.I2CAddressesRange(), sensors.LocateI2CSensor)
	for _, devices := range m.detectedI2Cs {
		for _, s := range devices {
			detectedSensors[s.ID()] = s
		}
	}

	for id := range registeredSensors {
		if !detectedSensors.Exists(id) && !m.contains(staticSensors, id) {
			payload.Removed = append(payload.Removed, id)
			isChanges = true
			shared.Logger.Debugf("Hotswap: %s sensor was detached from the device", id)
		}
	}

	for id := range detectedSensors {
		if !registeredSensors.Exists(id) {
			payload.Added = append(payload.Added, detectedSensors[id])
			isChanges = true
			shared.Logger.Debugf("Hotswap: %s sensor was attached to the device", id)
		}
	}

	if isChanges {
		eventdriver.EmitEvent(ctx, events.SensorsRegisterChanged, payload)
		m.UpdateSensorsRegister(payload.Added, payload.Removed)
	}

	return nil
}

func (m *HotswapDetector) contains(register map[string]sensor.Sensor, id string) bool {
	_, contains := register[id]
	return contains
}
