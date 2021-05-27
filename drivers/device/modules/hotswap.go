package modules

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	dev "github.com/timoth-y/chainmetric-sensorsys/drivers/device"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/periphery"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensors"
	"github.com/timoth-y/chainmetric-sensorsys/model"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// HotswapDetector defines device.Device module for detecting changes in connected sensors.
type HotswapDetector struct {
	dev  *dev.Device
	once *sync.Once

	detectedI2Cs  periphery.I2CDetectResults
}

// WithHotswapDetector can be used to setup HotswapDetector module for the device.Device.
func WithHotswapDetector() Module {
	return &HotswapDetector{
		once: &sync.Once{},
	}
}

func (m *HotswapDetector) Setup(device *dev.Device) error {
	m.dev = device

	return nil
}

func (m *HotswapDetector) Start(ctx context.Context) {
	m.once.Do(func() {
		var (
			startTime  time.Time
			interval = viper.GetDuration("device.hotswap_detect_interval")
		)

		go func() {
		LOOP:
			for {
				startTime = time.Now()

				if err := m.handleHotswap(); err != nil {
					shared.Logger.Error(errors.Wrap(err, "failed to handle hotswap"))
				}

				select {
				case <- time.After(interval - time.Since(startTime)):
				case <- ctx.Done():
					shared.Logger.Debug("Hotswap detector module routine ended.")
					break LOOP
				}
			}
		}()
	})
}

func (m *HotswapDetector) handleHotswap() error {
	var (
		detectedSensors = make(map[string]sensor.Sensor)
		registeredSensors = m.dev.RegisteredSensors()
		staticSensors = m.dev.StaticSensors()
		isChanges bool
	)

	m.detectedI2Cs = periphery.ScanI2C(sensors.I2CAddressesRange(), sensors.LocateI2CSensor)
	for _, devices := range m.detectedI2Cs {
		for _, s := range devices {
			detectedSensors[s.ID()] = s
		}
	}

	for id := range registeredSensors {
		if _, ok := detectedSensors[id]; !ok && !m.contains(staticSensors, id) {
			m.dev.UnregisterSensor(id)
			isChanges = true
			shared.Logger.Debugf("Hotswap: %s sensor was detached from the device", id)
		}
	}

	for id := range detectedSensors {
		if _, ok := registeredSensors[id]; !ok {
			m.dev.RegisterSensors(detectedSensors[id])
			isChanges = true
			shared.Logger.Debugf("Hotswap: %s sensor was attached to the device", id)
		}
	}

	if isChanges {
		if err := m.dev.SetSpecs(func(specs *model.DeviceSpecs) {
			specs.Supports = m.dev.RegisteredSensors().Union(staticSensors).SupportedMetrics()
		}); err != nil {
			return err
		}
	}

	return nil
}

func (m *HotswapDetector) contains(register map[string]sensor.Sensor, id string) bool {
	_, contains := register[id]
	return contains
}
