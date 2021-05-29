package device

import (
	"sync"
	"time"

	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-sensorsys/model"
)

type (
	// cacheLayer defines an extension layer for Device,
	// containing and providing methods to manage blockchain data cache.
	cacheLayer struct {
		assets   *assetsCache
		requests *requirementsCache
	}

	// assetsCache defines structure for storing models.Asset records data taken from blockchain.
	assetsCache struct {
		mutex sync.Mutex
		data  map[string]bool
	}

	// requirementsCache defines structure for storing models.Requirements records data taken from blockchain.
	requirementsCache struct {
		mutex sync.Mutex
		data  map[string]model.SensorsReadingRequest
	}
)

// newCacheLayer constructs new cacheLayer instance.
func newCacheLayer() cacheLayer {
	return cacheLayer{
		assets: &assetsCache{
			mutex: sync.Mutex{},
			data:  make(map[string]bool),
		},
		requests: &requirementsCache{
			mutex: sync.Mutex{},
			data:  make(map[string]model.SensorsReadingRequest),
		},
	}
}

// GetCachedAssets returns IDs of cached models.Asset records.
func (c *cacheLayer) GetCachedAssets() []string {
	c.assets.mutex.Lock()
	defer c.assets.mutex.Unlock()

	var (
		ids = make([]string, len(c.assets.data))
		i = 0
	)

	for id := range c.assets.data {
		ids[i] = id
		i++
	}

	return ids
}

// ExistsAssetInCache determines whether the models.Asset record stored in cache by given `id`.
func (c *cacheLayer) ExistsAssetInCache(id string) bool {
	c.assets.mutex.Lock()
	defer c.assets.mutex.Unlock()

	_, exists := c.assets.data[id]
	return exists
}

// PutAssetsToCache puts models.Asset records to local cache.
func (c *cacheLayer) PutAssetsToCache(assets ...*models.Asset) {
	c.assets.mutex.Lock()
	defer c.assets.mutex.Unlock()

	for i := range assets {
		c.assets.data[assets[i].ID] = true
	}
}

// RemoveAssetFromCache removes models.Asset record to local cache by given `id`.
func (c *cacheLayer) RemoveAssetFromCache(id string) {
	c.assets.mutex.Lock()
	defer c.assets.mutex.Unlock()

	delete(c.assets.data, id)
}

// FlushAssetsCache resets models.Asset records cache.
func (c *cacheLayer) FlushAssetsCache() {
	c.assets.mutex.Lock()
	defer c.assets.mutex.Unlock()

	c.assets.data = make(map[string]bool)
}

// GetCachedRequirements returns data of cached models.Requirements records as model.SensorsReadingRequest.
func (c *cacheLayer) GetCachedRequirements() []model.SensorsReadingRequest {
	c.requests.mutex.Lock()
	defer c.requests.mutex.Unlock()

	var (
		reqs = make([]model.SensorsReadingRequest, len(c.requests.data))
		i    = 0
	)

	for _, req := range c.requests.data {
		reqs[i] = req
		i++
	}

	return reqs
}

// GetRequirementsFromCache tries to retrieve single models.Requirements from cache by given `id`,
// where it is stored as a model.SensorsReadingRequest record.
func (c *cacheLayer) GetRequirementsFromCache(id string) (model.SensorsReadingRequest, bool) {
	c.requests.mutex.Lock()
	defer c.requests.mutex.Unlock()

	req, exists := c.requests.data[id]
	return req, exists
}

// PutRequirementsToCache puts models.Requirements records to local cache
// and return requests as they are stored.
func (c *cacheLayer) PutRequirementsToCache(reqs ...*models.Requirements) (requests []model.SensorsReadingRequest) {
	c.requests.mutex.Lock()
	defer c.requests.mutex.Unlock()

	for _, req := range reqs {
		request := model.SensorsReadingRequest{
			AssetID: req.AssetID,
			Metrics: req.Metrics.Metrics(),
			Period:  time.Second * time.Duration(req.Period),
		}

		c.requests.data[req.ID] = request
		requests = append(requests, request)
	}

	return requests
}

// RemoveRequirementsFromCache removes models.Requirements record to local cache by given `id`.
func (c *cacheLayer) RemoveRequirementsFromCache(id string) {
	c.requests.mutex.Lock()
	defer c.requests.mutex.Unlock()

	delete(c.requests.data, id)
}

// FlushRequirementsCache resets models.Requirements records cache.
func (c *cacheLayer) FlushRequirementsCache() {
	c.requests.mutex.Lock()
	defer c.requests.mutex.Unlock()

	c.requests.data = make(map[string]model.SensorsReadingRequest)
}
