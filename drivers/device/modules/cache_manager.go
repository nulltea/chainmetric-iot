package modules

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models/requests"
	dev "github.com/timoth-y/chainmetric-sensorsys/drivers/device"
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
func WithCacheManager() dev.Module {
	return &CacheManager{
		moduleBase: withModuleBase("cache_manager"),
	}
}

func (m *CacheManager) Start(ctx context.Context) {
	go m.Do(func() {
		if !m.trySyncWithDeviceLifecycle(ctx, m.Start) {
			return
		}

		m.cacheBlockchainData()
		eventdriver.SubscribeHandler(events.DeviceLocationChanged, m.handleReCachingEvents)
	})
}

func (m *CacheManager) cacheBlockchainData() {
	if err := m.locateAssets(); err != nil {
		shared.Logger.Error(errors.Wrap(err, "failed to cache assets"))
	}

	if err := m.receiveRequirements(); err != nil {
		shared.Logger.Error(errors.Wrap(err, "failed to cache requirements"))
	}
}

func (m *CacheManager) locateAssets() error {
	var (
		contract = blockchain.Contracts.Assets
	)

	m.FlushAssetsCache()

	assets, err := contract.Receive(requests.AssetsQuery{
		Location: &requests.LocationQuery{
			GeoPoint: m.Location(),
			Distance: viper.GetFloat64("assets_locate_distance"),
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

	reqs, err := contract.ReceiveFor(m.GetCachedAssets())
	if err != nil {
		return err
	}

	m.PutRequirementsToCache(reqs...)

	shared.Logger.Debugf("Received %d requirements", len(reqs))

	return nil
}

func (m *CacheManager) handleReCachingEvents(_ context.Context, v interface{}) error {
	if _, ok := v.(events.DeviceLocationChangedPayload); ok {
		m.cacheBlockchainData()
		return nil
	}

	return eventdriver.ErrIncorrectPayload
}
