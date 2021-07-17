package modules

import (
	"context"
	"fmt"
	"strings"
	"sync"
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
	viewLock *sync.Mutex

	requestsThroughput []float64
}

// WithGUIRenderer can be used to setup GUIRenderer logical device.Module onto the device.Device.
func WithGUIRenderer() device.Module {
	return &GUIRenderer{
		moduleBase: withModuleBase("GUI_RENDERER"),
		viewLock: &sync.Mutex{},
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

		// Act on changes sensors pool to view hotswap notification:
		eventdriver.SubscribeHandler(events.SensorsRegisterChanged, func(_ context.Context, v interface{}) error {
			if payload, ok := v.(events.SensorsRegisterChangedPayload); ok {
				m.renderHotswapNotification(payload)
				return nil
			}

			return eventdriver.ErrIncorrectPayload
		})

		m.renderStats()
		m.renderLoop(ctx)
	})
}

func (m *GUIRenderer) renderLoop(ctx context.Context) {
	var (
		ticker = time.NewTicker(viper.GetDuration("device.gui_update_interval"))
	)

LOOP:
	for {
		select {
		case <- ticker.C:
			m.renderStats()
		case <- ctx.Done():
			shared.Logger.Debug("GUI renderer module routine ended")
			break LOOP
		}
	}
}

func (m *GUIRenderer) renderStats() {
	var (
		builder  = strings.Builder{}
		interval = viper.GetDuration("device.gui_update_interval")
		throughput []float64
	)

	m.viewLock.Lock()
	defer m.viewLock.Unlock()

	for _, count := range m.requestsThroughput {
		throughput = append(throughput, count * 60 / interval.Seconds())
	}

	builder.WriteString(fmt.Sprintf("IP: %s\n", m.Specs().IPAddress))
	builder.WriteString(fmt.Sprintf("Supported: %d metrics\n", len(m.Specs().Supports)))
	builder.WriteString(fmt.Sprintf("Thoughput: %d requests\\min",
		int(throughput[len(m.requestsThroughput) - 1]),
	))

	gui.RenderWithChart(builder.String(), m.requestsThroughput...)

	m.requestsThroughput = append(m.requestsThroughput, 0)
}

func (m *GUIRenderer) renderHotswapNotification(event events.SensorsRegisterChangedPayload) {
	var (
		builder = strings.Builder{}
		attached []string
	)

	m.viewLock.Lock()
	defer m.viewLock.Unlock()

	for i := range event.Added {
		attached = append(attached, event.Added[i].ID())
	}

	if len(attached) > 0 {
		var (
			word = "have"
		)

		if len(attached) > 1 {
			word = "have"
		}

		builder.WriteString(fmt.Sprintf("%s %s just been attached",
			strings.Join(attached, ","),
			word,
		))
	}

	if len(event.Removed) > 0 {
		var (
			word = "was"
		)

		if len(attached) > 1 {
			word = "were"
		}

		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("%s %s removed",
			strings.Join(event.Removed, ","),
			word,
		))
	}

	gui.RenderTextWithIcon(builder.String(), "hotswap")
}
