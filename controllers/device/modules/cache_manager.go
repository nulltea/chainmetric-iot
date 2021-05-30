package modules

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models/requests"
	"github.com/timoth-y/chainmetric-sensorsys/controllers/device"
	"github.com/timoth-y/chainmetric-sensorsys/model/events"
	"github.com/timoth-y/chainmetric-sensorsys/network/blockchain"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
	"github.com/timoth-y/go-eventdriver"
)

// CacheManager implements device.Module for device.Device cache data managing.
type CacheManager struct {
	moduleBase
}

// WithCacheManager can be used to setup CacheManager logical device.Module onto the device.Device.
func WithCacheManager() device.Module {
	return &CacheManager{
		moduleBase: withModuleBase("CACHE_MANAGER"),
	}
}

func (m *CacheManager) Start(ctx context.Context) {
	go m.Do(func() {
		if !m.trySyncWithDeviceLifecycle(ctx, m.Start) {
			return
		}

		// Handle changes that require full cache reload:
		eventdriver.SubscribeHandler(events.DeviceLocationChanged, func(_ context.Context, _ interface{}) error {
			// Canceling requests before re-caching:
			for _, request := range m.GetCachedRequirementsFor() {
				request.Cancel()
			}

			m.cacheBlockchainData(ctx)
			return nil
		})

		// Handle changes in assigned assets, which require changes in requirements:
		eventdriver.SubscribeHandler(events.AssetsChanged, func(ctx context.Context, v interface{}) error {
			if payload, ok := v.(events.AssetsChangedPayload); ok {
				// Canceling requests for removed assets and removing them from cache:
				for _, request := range m.GetCachedRequirementsFor(payload.Removed...) {
					request.Cancel()
					m.RemoveRequirementsFromCache(request.ID)
				}

				res, err := blockchain.Contracts.Requirements.ReceiveFor(payload.Assigned...)
				if err != nil {
					return errors.Wrap(err, "failed to receive requirements for newly assigned assets")
				}

				eventdriver.EmitEvent(ctx, events.RequirementsChanged, events.RequirementsChangedPayload{
					Requests: m.PutRequirementsToCache(res...),
				})

				return nil
			}

			return eventdriver.ErrIncorrectPayload
		})

		m.cacheBlockchainData(ctx)
	})
}

func (m *CacheManager) cacheBlockchainData(ctx context.Context) {
	defer eventdriver.EmitEvent(ctx, events.CacheChanged, nil)

	if err := m.locateAssets(); err != nil {
		shared.Logger.Error(errors.Wrap(err, "failed to cache assets"))
	}

	if err := m.receiveRequirements(); err != nil {
		shared.Logger.Error(errors.Wrap(err, "failed to cache requirements"))
	}

	eventdriver.EmitEvent(ctx, events.CacheChanged, nil)
}

func (m *CacheManager) locateAssets() error {
	var (
		contract = blockchain.Contracts.Assets
	)

	m.FlushAssetsCache()

	assets, err := contract.Receive(requests.AssetsQuery{
		Location: &requests.LocationQuery{
			GeoPoint: m.Location(),
			Distance: viper.GetFloat64("device.assets_locate_distance"),
		},
	})
	if err != nil {
		return err
	}

	m.PutAssetsToCache(assets...)

	shared.Logger.Debugf("Located %d assets on location: %s", len(assets), m.Location().Name)

	return nil
}

func (m *CacheManager) receiveRequirements() error {
	var (
		contract = blockchain.Contracts.Requirements
	)

	m.FlushRequirementsCache()

	reqs, err := contract.ReceiveFor(m.GetCachedAssets()...)
	if err != nil {
		return err
	}

	m.PutRequirementsToCache(reqs...)

	shared.Logger.Debugf("Received %d requirements", len(reqs))

	return nil
}
