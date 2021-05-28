package modules

import (
	"context"
	"sync"
	"time"

	fabricStatus "github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/utils"
	dev "github.com/timoth-y/chainmetric-sensorsys/drivers/device"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/storage"
	"github.com/timoth-y/chainmetric-sensorsys/model/events"
	"github.com/timoth-y/chainmetric-sensorsys/network/blockchain"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
	"github.com/timoth-y/go-eventdriver"
)

// FailoverHandler defines device.Device module for handling operating failures.
type FailoverHandler struct {
	*dev.Device
	once      *sync.Once
	ctx       context.Context
	pingTimer *time.Timer
}

// WithFailoverHandler can be used to setup FailoverHandler module for the device.Device.
func WithFailoverHandler() Module {
	return &FailoverHandler{
		once: &sync.Once{},
		ctx: context.Background(),
	}
}

func (m *FailoverHandler) Setup(device *dev.Device) error {
	m.Device = device

	return nil
}

func (m *FailoverHandler) Start(ctx context.Context) {
	m.once.Do(func() {
		m.ctx = ctx

		// Try post leftover readings in cache
		go m.tryRepostCachedReadings()

		// Listen to metric readings failures
		eventdriver.SubscribeHandler(events.MetricReadingsPostFailed, func(ctx context.Context, v interface{}) error {
			if payload, ok := v.(events.MetricReadingsPostFailedPayload); ok {
				m.handleFailedToPostReadings(payload.MetricReadings)
				return nil
			}

			return eventdriver.ErrIncorrectPayload
		})
	})
}


func (m *FailoverHandler) handleFailedToPostReadings(readings models.MetricReadings) {
	m.pingNetworkConnection()

	if err := storage.CacheReadings(readings); err != nil {
		shared.Logger.Error(errors.Wrap(err, "failed to cache readings while network connection absence"))
		return
	}

	shared.Logger.Warningf(
		"Detected network connection absence, cached readings for %s to post later => %s", readings.AssetID,
		utils.Prettify(readings),
	)
}

func (m *FailoverHandler) pingNetworkConnection() {
	var (
		interval = viper.GetDuration("device.ping_timer_interval")
	)

	if m.pingTimer != nil {
		if !m.pingTimer.Reset(interval) {
			go m.ping(m.pingTimer, m.tryRepostCachedReadings)
		}
	} else {
		m.pingTimer = time.NewTimer(interval)
		go m.ping(m.pingTimer, m.tryRepostCachedReadings)
	}
}

// tryRepostCachedReadings makes attempt to repost cached during network absence sensor readings data.
func (m *FailoverHandler) tryRepostCachedReadings() {
	storage.IterateOverCachedReadings(m.ctx, func(key string, record models.MetricReadings) (toBreak bool, err error) {
		if err = blockchain.Contracts.Readings.Post(record); err != nil {
			if detectNetworkAbsence(err) {
				m.pingNetworkConnection()
				shared.Logger.Debug("Network connection is still down - stop iterating sequence")

				return true, nil
			}

			return false, err
		}

		shared.Logger.Debugf("Successfully posted cached readings for key: %s => %s", key, utils.Prettify(record))

		return false, nil
	}, true)
}

func (m *FailoverHandler) ping(t *time.Timer, onPong func()) {
	<- t.C
	onPong()
}

func detectNetworkAbsence(err error) bool {
	if status, ok := fabricStatus.FromError(err); ok {
		switch status.Group {
		case fabricStatus.DiscoveryServerStatus:
			return true
		}
	}

	return false
}
