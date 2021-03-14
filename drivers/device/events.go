package device

import (
	"context"
	"sync"

	"github.com/timoth-y/iot-blockchain-contracts/models"
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
		go d.subscribeToAssets(ctx)
	})
}

func (d *Device) subscribeToAssets(ctx context.Context) {
	var (
		contract = d.client.Contracts.Assets
	)

	contract.Subscribe(ctx, "*", func(asset *models.Asset, e string) error {
		switch e {
		case "inserted":
			fallthrough
		case "updated":
			if asset.Location == d.model.Location {
				d.assets[asset.ID] = true
				break
			}
			fallthrough
		case "removed":
			delete(d.assets, asset.ID)
		}

		return nil
	})
}
