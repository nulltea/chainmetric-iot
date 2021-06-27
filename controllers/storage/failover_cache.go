package storage

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/utils"

	"github.com/timoth-y/chainmetric-iot/shared"
)

// ReadingsCacheIteratorFunc defines function called by sensor readings cache iterator.
type ReadingsCacheIteratorFunc func(key string, record models.MetricReadings) (toBreak bool, err error)

// CacheReadings stores models.MetricReadings into local cache DB.
func CacheReadings(readings ...models.MetricReadings) (err error) {
	var (
		batch = new(leveldb.Batch)
	)

	for _, reading := range readings {
		var (
			key = utils.FormCompositeKey("reading",
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

// IterateOverCachedReadings performs iteration over all cached models.MetricReadings records.
// allowing to `pop` them on fly.
func IterateOverCachedReadings(ctx context.Context, fn ReadingsCacheIteratorFunc, pop bool) {
	var (
		prefix = []byte(utils.FormCompositeKey("reading"))
		iter = shared.LevelDB.NewIterator(util.BytesPrefix(prefix), nil)
	)

	for iter.Next() {
		select {
		case <- ctx.Done():
			return
		default:
		}

		var (
			key = string(iter.Key())
			_, attrs = utils.SplitCompositeKey(key)
			values map[models.Metric]float64
		)

		if len(attrs) < 2 {
			shared.Logger.Warningf("Invalid composite key '%s': it must contain assetID and timestamp", key)
			continue
		}

		var (
			assetID = attrs[0]
			unix, _ = strconv.Atoi(attrs[1])
			timestamp = time.Unix(int64(unix), 0)
		)

		if len(assetID) == 0 || timestamp.IsZero() {
			shared.Logger.Warningf("Invalid composite key data in '%s'", key)
			continue
		}

		if err := json.Unmarshal(iter.Value(), &values); err != nil {
			shared.Logger.Error(errors.Wrapf(err, "failed to unmarshal values for key '%s'", key))
			continue
		}

		toBreak, err := fn(key, models.MetricReadings{
			AssetID: assetID,
			Timestamp: timestamp,
			Values: values,
		})

		if err != nil {
			shared.Logger.Error(errors.Wrapf(err,
				"something went wrong when iterating over '%s' - stop iterating sequence", key))
		}

		if toBreak {
			return
		}

		if pop {
			if err := shared.LevelDB.Delete([]byte(key), nil); err != nil {
				shared.Logger.Error(errors.Wrapf(err,
					"failed to delete on key '%s' - stop iterating sequence", key))
				return
			}
		}
	}
}
