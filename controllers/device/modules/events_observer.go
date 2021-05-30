package modules

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-sensorsys/controllers/device"
	"github.com/timoth-y/chainmetric-sensorsys/model/events"
	"github.com/timoth-y/chainmetric-sensorsys/network/blockchain"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
	"github.com/timoth-y/go-eventdriver"
)

// EventsObserver implements device.Module for listening and acting on changes in blockchain ledger data.
//
// This device.Module also capable of mutating cache layer data of the device.Device.
type EventsObserver struct {
	moduleBase
}

// WithEventsObserver can be used to setup EventsObserver logical device.Module onto the device.Device.
func WithEventsObserver() device.Module {
	return &EventsObserver{
		moduleBase: withModuleBase("EVENTS_OBSERVER"),
	}
}

func (m *EventsObserver) Start(ctx context.Context) {
	go m.Do(func() {
		if !m.trySyncWithDeviceLifecycle(ctx, m.Start) {
			return
		}

		go m.watchAssets(ctx)
		go m.watchDevice(ctx)
		go m.watchRequirements(ctx)
	})
}

func (m *EventsObserver) watchAssets(ctx context.Context) {
	blockchain.Contracts.Assets.Subscribe(ctx, "*", func(asset *models.Asset, e string) error {
		switch e {
		case "inserted":
			fallthrough
		case "updated":
			if asset.Location.IsNearBy(m.Location(), viper.GetFloat64("assets_locate_distance")) {
				m.PutAssetsToCache(asset)
				break
			}
			fallthrough
		case "removed":
			m.RemoveAssetFromCache(asset.ID)
		}

		shared.Logger.Debugf("Asset %q was %s", asset.ID, e)

		return nil
	})
}

func (m *EventsObserver) watchDevice(ctx context.Context) {
	if err := blockchain.Contracts.Devices.Subscribe(ctx, "*", func(dev *models.Device, e string) error {
		if dev.ID != m.ID() {
			return nil
		}

		switch e {
		case "updated":
			m.actOnDeviceUpdates(ctx, dev)
			fallthrough
		case "inserted":
			m.UpdateDeviceModel(dev)
		case "removed":
			shared.Logger.Notice("Device has been removed from blockchain, must reset it now")
			eventdriver.EmitEvent(ctx, events.DeviceRemovedFromNetwork, nil)
		}

		shared.Logger.Debugf("Device was %s", e)

		return nil
	}); err != nil {
		shared.Logger.Fatal(errors.Wrap(err, "failed to subscribe to device changes on network"))
	}
}

func (m *EventsObserver) watchRequirements(ctx context.Context) {
	blockchain.Contracts.Requirements.Subscribe(ctx, "*", func(req *models.Requirements, e string) error {
		if !m.ExistsAssetInCache(req.AssetID) {
			return nil
		}

		switch e {
		case "updated":
			if request, ok := m.GetRequirementsFromCache(req.ID); ok {
				request.Cancel()
			}
			fallthrough
		case "inserted":
			eventdriver.EmitEvent(ctx, events.RequirementsSubmitted, events.RequirementsSubmittedPayload{
				Requests: m.PutRequirementsToCache(req),
			})
			shared.Logger.Debugf(
				"Requirements (id: %s) with %d metrics was %s", req.ID, len(req.Metrics), e,
			)
		case "removed":
			if request, ok := m.GetRequirementsFromCache(req.ID); ok {
				request.Cancel()
				m.RemoveRequirementsFromCache(req.ID)
				shared.Logger.Debugf(
					"Requirements (id: %s) was removed and unsubscribed from reading sensors", req.ID,
				)
			}
		}

		return nil
	})
}

func (m *EventsObserver) actOnDeviceUpdates(ctx context.Context, updated *models.Device) {
	if m.Location().IsNearBy(updated.Location, viper.GetFloat64("assets_locate_distance")) {
		eventdriver.EmitEvent(ctx, events.DeviceLocationChanged, events.DeviceLocationChangedPayload{
			Old: m.Location(),
			New: updated.Location,
		})
	}
}
