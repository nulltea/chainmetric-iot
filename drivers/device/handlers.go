package device

import (
	"time"

	fabricStatus "github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/storage"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

func (d *Device) handleNetworkDisconnection(readings models.MetricReadings) {
	d.pingNetworkConnection()

	if err := storage.CacheReadings(readings); err != nil {
		shared.Logger.Error(errors.Wrap(err, "failed to cache readings while network connection absence"))
		return
	}

	shared.Logger.Warningf(
		"Detected network connection absence, cached readings for %s to post later => %s", readings.AssetID,
		shared.Prettify(readings),
	)
}

func (d *Device) pingNetworkConnection() {
	var (
		interval = viper.GetDuration("device.ping_timer_interval")
	)

	if d.pingTimer != nil {
		if !d.pingTimer.Reset(interval) {
			go ping(d.pingTimer, d.tryRepostCachedReadings)
		}
	} else {
		d.pingTimer = time.NewTimer(interval)
		go ping(d.pingTimer, d.tryRepostCachedReadings)
	}
}

// tryRepostCachedReadings makes attempt to repost cached during network absence sensor readings data.
func (d *Device) tryRepostCachedReadings() {
	var (
		contract = d.client.Contracts.Readings
	)

	storage.IterateOverCachedReadings(func(key string, record models.MetricReadings) (toBreak bool, err error) {
		if err = contract.Post(record); err != nil {
			if detectNetworkAbsence(err) {
				d.pingNetworkConnection()
				shared.Logger.Debug("Network connection is still down - stop iterating sequence")

				return true, nil
			}

			return false, err
		}

		shared.Logger.Debugf("Successfully posted cached readings for key: %s => %s", key, shared.Prettify(record))

		return false, nil
	}, true)
}

func detectNetworkAbsence(err error) bool {
	if status, ok := fabricStatus.FromError(err); ok {
		switch status.Group {
		case fabricStatus.DiscoveryServerStatus:
			return true
		}
	}

	return false
}

func ping(t *time.Timer, onPong func()) {
	<- t.C
	onPong()
}
