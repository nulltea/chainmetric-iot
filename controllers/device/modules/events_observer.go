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
	if err := blockchain.Contracts.Assets.Subscribe(ctx, "*", func(asset *models.Asset, e string) error {
		var (
			changesPayload = events.AssetsChangedPayload{}
		)

		switch e {
		case "inserted", "updated":
			if asset.Location.IsNearBy(m.Location(), viper.GetFloat64("device.assets_locate_distance")) {
				if !m.ExistsAssetInCache(asset.ID) {
					changesPayload.Assigned = append(changesPayload.Assigned, asset.ID)
					shared.Logger.Debugf("Asset %q was assigned for the device", asset.ID)
				} else {
					shared.Logger.Debugf("Already assigned asset %q was changed ", asset.ID)
				}

				m.PutAssetsToCache(asset)
				break
			}
			fallthrough
		case "removed":
			if m.ExistsAssetInCache(asset.ID) {
				m.RemoveAssetFromCache(asset.ID)
				changesPayload.Removed = append(changesPayload.Removed, asset.ID)
				shared.Logger.Debugf("Asset %q was unassigned from the device", asset.ID)
			}
		}

		eventdriver.EmitEvent(ctx, events.AssetsChanged, changesPayload)
		return nil
	}); err != nil {
		shared.Logger.Error(errors.Wrap(err, "failed to subscribe to assets changes on network"))
	}
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
	if err := blockchain.Contracts.Requirements.Subscribe(ctx, "*",
		func(req *models.Requirements, e string) error {
			switch e {
			case "inserted", "updated":
				if !m.ExistsAssetInCache(req.AssetID) {
					break
				}

				// Canceling already existing requests since they will be handled again letter on:
				if request, ok := m.GetRequirementsFromCache(req.ID); ok {
					request.Cancel()
				}

				// Putting new or changed requirements to cache and notifying other modules about changes:
				eventdriver.EmitEvent(ctx, events.RequirementsChanged, events.RequirementsChangedPayload{
					Requests: m.PutRequirementsToCache(req),
				})

				shared.Logger.Debugf("Requirements (id: %s) with %d metrics was %s", req.ID,
					len(req.Metrics), e,
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
		}); err != nil {
		shared.Logger.Fatal(errors.Wrap(err, "failed to subscribe to requirements changes on network"))
	}
}

func (m *EventsObserver) actOnDeviceUpdates(ctx context.Context, updated *models.Device) {
	if !m.Location().IsNearBy(updated.Location, viper.GetFloat64("device.assets_locate_distance")) {
		eventdriver.EmitEvent(ctx, events.DeviceLocationChanged, events.DeviceLocationChangedPayload{
			Old: m.Location(),
			New: updated.Location,
		})
	}
}
