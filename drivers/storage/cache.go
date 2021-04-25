package storage

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

// CacheReadings stores models.MetricReadings into local cache DB
func CacheReadings(readings ...models.MetricReadings) (err error) {
	var (
		batch = new(leveldb.Batch)
	)

	for _, reading := range readings {
		var (
			key = shared.FormCompositeKey("reading",
				reading.AssetID,
				strconv.Itoa(int(reading.Timestamp.Unix())),
			)
			value []byte
		)

		if value, err = json.Marshal(reading.Values); err != nil {
			return err
		}

		batch.Put([]byte(key), value)
	}

	return shared.LevelDB.Write(batch, nil)
}

func IterateOverCachedReadings(fn func(models.MetricReadings) error) {
	var (
		prefix = []byte(shared.FormCompositeKey("reading"))
		iter = shared.LevelDB.NewIterator(util.BytesPrefix(prefix), nil)
	)

	for iter.Next() {
		var (
			key = string(iter.Key())
			_, attrs = shared.SplitCompositeKey(key)
			values map[models.Metric]interface{}
		)

		if len(attrs) < 2 {
			shared.Logger.Warningf("Invalid composite key %q: it must contain assetID and timestamp", key)
			continue
		}

		var (
			assetID = attrs[0]
			unix, _ = strconv.Atoi(attrs[1])
			timestamp = time.Unix(int64(unix), 0)
		)

		if len(assetID) == 0 || timestamp.IsZero() {
			shared.Logger.Warningf("Invalid composite key data in %q", key)
			continue
		}

		if err := json.Unmarshal(iter.Value(), &values); err != nil {
			shared.Logger.Error(errors.Wrapf(err, "failed to unmarshal values for key %q", key))
			continue
		}

		if err := fn(models.MetricReadings{
			AssetID: assetID,
			Timestamp: timestamp,
			Values: values,
		}); err != nil {
			shared.Logger.Error(errors.Wrapf(err, "something went wrong when iterating over %q", key))
		}
	}
}
