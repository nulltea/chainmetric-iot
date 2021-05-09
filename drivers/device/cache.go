package device

import (
	"context"
	"sync"
	"time"

	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/models/requests"

	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

type assetsCache struct {
	mutex sync.Mutex
	data  map[string]bool
}

type requirementsCache struct {
	mutex sync.Mutex
	data  map[string]*readingsRequest
}

type readingsRequest struct {
	id      string
	assetID string
	period  time.Duration
	metrics models.Metrics
	cancel  context.CancelFunc
}

func (d *Device) CacheBlockchainState() error {
	if !d.active {
		return nil
	}

	if err := d.locateAssets(); err != nil {
		return err
	}

	if err := d.receiveRequirements(); err != nil {
		return err
	}

	return nil
}

func (d *Device) locateAssets() error {
	var (
		contract = d.client.Contracts.Assets
	)
	d.assets.mutex.Lock()
	defer d.assets.mutex.Unlock()

	d.assets.data = make(map[string]bool)

	assets, err := contract.Receive(requests.AssetsQuery{
		Location: &d.model.Location.Name,
	}); if err != nil {
		return err
	}

	for _, asset := range assets {
		d.assets.data[asset.ID] = true
	}

	shared.Logger.Debugf("Located %d assets on location: %s", len(assets), d.model.Location.Name)

	return nil
}

func (d *Device) receiveRequirements() error {
	var (
		contract = d.client.Contracts.Requirements
	)
	d.requests.mutex.Lock()
	defer d.requests.mutex.Unlock()

	d.requests.data = make(map[string]*readingsRequest)

	reqs, err := contract.ReceiveFor(d.assets.Get()); if err != nil {
		return err
	}

	for _, req := range reqs {
		d.requests.data[req.ID] = &readingsRequest{
			assetID: req.AssetID,
			metrics: req.Metrics.Metrics(),
			period: time.Second * time.Duration(req.Period),
		}
	}

	shared.Logger.Debugf("Received %d requirements", len(reqs))

	return nil
}


func (ac *assetsCache) Get() []string {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()

	var (
		ids = make([]string, len(ac.data))
		i = 0
	)

	for id := range ac.data {
		ids[i] = id
		i++
	}
	return ids
}

func (rc *requirementsCache) Get() []*readingsRequest {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	var (
		reqs = make([]*readingsRequest, len(rc.data))
		i    = 0
	)

	for _, req := range rc.data {
		reqs[i] = req
		i++
	}
	return reqs
}

