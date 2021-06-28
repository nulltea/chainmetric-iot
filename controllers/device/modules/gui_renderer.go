package modules

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-iot/controllers/device"
	"github.com/timoth-y/chainmetric-iot/controllers/gui"
	"github.com/timoth-y/chainmetric-iot/model/events"
	"github.com/timoth-y/chainmetric-iot/shared"
	"github.com/timoth-y/go-eventdriver"
)

// GUIRenderer implements device.Module for device.Device GUI controlling.
type GUIRenderer struct {
	moduleBase

	requestsThroughput []float64
}

// WithGUIRenderer can be used to setup GUIRenderer logical device.Module onto the device.Device.
func WithGUIRenderer() device.Module {
	return &GUIRenderer{
		moduleBase: withModuleBase("GUI_RENDERER"),
	}
}

func (m *GUIRenderer) Setup(device *device.Device) error {
	if !gui.Available() {
		return errors.New("module won't work without display available")
	}

	if err := m.moduleBase.Setup(device); err != nil {
		return err
	}

	m.requestsThroughput = append(m.requestsThroughput, 0)

	return nil
}

func (m *GUIRenderer) Start(ctx context.Context) {
	go m.Do(func() {
		if !m.trySyncWithDeviceLifecycle(ctx, m.Start) {
			return
		}

		// Act on each new handled request to update device throughput:
		eventdriver.SubscribeHandler(events.RequestHandled, func(_ context.Context, _ interface{}) error {
			m.requestsThroughput[len(m.requestsThroughput) - 1]++
			return nil
		})

		m.renderLoop(ctx)
	})
}

func (m *GUIRenderer) renderLoop(ctx context.Context) {
	var (
		interval = viper.GetDuration("device.gui_update_interval")
		startTime  time.Time
	)

LOOP:
	for {
		startTime = time.Now()

		gui.RenderWithChart(fmt.Sprintf(
`IP: %s
Supported: %d metrics
Thoughput: %d requests\min`,
				m.Specs().IPAddress,
				len(m.Specs().Supports),
				int(m.requestsThroughput[len(m.requestsThroughput) - 1] * 60 / interval.Seconds()),
			), m.requestsThroughput...,
		)

		m.requestsThroughput = append(m.requestsThroughput, 0)

		select {
		case <- time.After(interval - time.Since(startTime)):
		case <- ctx.Done():
			shared.Logger.Debug("GUI renderer module routine ended")
			break LOOP
		}
	}
}
