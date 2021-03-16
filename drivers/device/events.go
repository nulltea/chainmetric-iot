package device

import (
	"context"
	"sync"
	"time"

	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

var (
	once = sync.Once{}
)

func (d *Device) WatchForBlockchainEvents() {
	var (
		ctx context.Context
	)

	ctx, d.cancelEvents = context.WithCancel(context.Background())

	once.Do(func() {
		go d.watchAssets(ctx)
		go d.watchDevice(ctx)
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

		return nil
	})
}

func (d *Device) watchDevice(ctx context.Context) {
	var (
		contract = d.client.Contracts.Devices
	)

	contract.Subscribe(ctx, "*", func(dev *models.Device, e string) error {
		switch e {
		case "inserted":
		case "updated":
			d.model = dev
		case "removed":
			shared.Logger.Notice("Device has been removed from blockchain, must reset it now")
			d.Reset()
			d.Close()
		}

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
		case "removed":
			if request, ok := d.requests.data[req.ID]; ok {
				request.cancel()
				delete(d.requests.data, req.ID)
			}
		}

		return nil
	})
}
