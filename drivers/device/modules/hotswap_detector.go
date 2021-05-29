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
	"github.com/timoth-y/chainmetric-sensorsys/model/events"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
	"github.com/timoth-y/go-eventdriver"
)

// HotswapDetector implements Module ofr detecting changes in connected to device.Device sensors.
type HotswapDetector struct {
	*dev.Device
	*sync.Once

	detectedI2Cs  periphery.I2CDetectResults
}

// WithHotswapDetector can be used to setup HotswapDetector logical Module onto the device.Device.
func WithHotswapDetector() Module {
	return &HotswapDetector{
		Once: &sync.Once{},
	}
}

func (m *HotswapDetector) MID() string {
	return "hotswap_detector"
}

func (m *HotswapDetector) Setup(device *dev.Device) error {
	m.Device = device

	return nil
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
				shared.Logger.Debug("Hotswap detector module routine ended.")
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

	m.detectedI2Cs = periphery.ScanI2C(sensors.I2CAddressesRange(), sensors.LocateI2CSensor)
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

// waitUntilSensorsDetected checks whether the sensors detected with a specific intervals.
func waitUntilSensorsDetected(d *dev.Device) bool {
	var attempts = 5
	for attempts > 0 {
		if d.RegisteredSensors().NotEmpty() {
			return true
		}

		time.Sleep(250 * time.Millisecond)
		attempts--
	}

	return false
}
