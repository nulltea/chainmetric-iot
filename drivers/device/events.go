package device

import (
	"context"
	"sync"
	"time"

	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

var (
	eventsOnce = sync.Once{}
)

func (d *Device) WatchForBlockchainEvents() {
	if !d.active {
		return
	}

	eventsOnce.Do(func() {
		go d.watchAssets(d.ctx)
		go d.watchDevice(d.ctx)
		go d.watchRequirements(d.ctx)
	})
}

func (d *Device) watchAssets(ctx context.Context) {
	var (
		contract = d.client.Contracts.Assets
	)

	contract.Subscribe(ctx, "*", func(asset *models.Asset, e string) error {
		d.assets.mutex.Lock()
		defer d.assets.mutex.Unlock()

		switch e {
		case "inserted":
			fallthrough
		case "updated":
			if asset.Location == d.model.Location {
				d.assets.data[asset.ID] = true
				break
			}
			fallthrough
		case "removed":
			delete(d.assets.data, asset.ID)
		}

		shared.Logger.Debugf("Asset %q was %s", asset.ID, e)

		return nil
	})
}

func (d *Device) watchDevice(ctx context.Context) {
	var (
		contract = d.client.Contracts.Devices
	)

	contract.Subscribe(ctx, "*", func(dev *models.Device, e string) error {
		switch e {
		case "updated":
			d.actOnDeviceUpdates(dev)
			fallthrough
		case "inserted":
			d.model = dev
		case "removed":
			shared.Logger.Notice("Device has been removed from blockchain, must reset it now")
			d.Reset()
			d.Close()
		}

		shared.Logger.Debugf("Device was %s", e)

		return nil
	})
}

func (d *Device) watchRequirements(ctx context.Context) {
	var (
		contract = d.client.Contracts.Requirements
	)

	contract.Subscribe(ctx, "*", func(req *models.Requirements, e string) error {
		d.requests.mutex.Lock()
		defer d.requests.mutex.Unlock()

		switch e {
		case "updated":
			if request, ok := d.requests.data[req.ID]; ok {
				request.cancel()
				delete(d.requests.data, req.ID)
			}
			fallthrough
		case "inserted":
			request := &readingsRequest{
				assetID: req.AssetID,
				metrics: req.Metrics.Metrics(),
				period: time.Second * time.Duration(req.Period),
			}
			d.requests.data[req.ID] = request
			d.actOnRequest(request)
			shared.Logger.Debugf("Requirements (id: %s) with %d metrics was %s", req.ID, len(req.Metrics), e)
		case "removed":
			if request, ok := d.requests.data[req.ID]; ok {
				request.cancel()
				delete(d.requests.data, req.ID)
			}
			shared.Logger.Debugf("Requirements (id: %s) was removed and unsubscribed from reading sensors",
				req.ID)
		}

		return nil
	})
}

func (d *Device) actOnDeviceUpdates(updated *models.Device) {
	if d.model.State != updated.State {
		// TODO: handle state changes
	}

	if d.model.Location != updated.Location {
		d.reader.Close() // TODO: main routine must stay locked from ending
		d.locateAssets()
		d.receiveRequirements()
		go d.Operate()
	}
}
