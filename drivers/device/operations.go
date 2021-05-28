package device

import (
	"time"

	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/utils"
	"github.com/timoth-y/chainmetric-sensorsys/network/blockchain"

	"github.com/timoth-y/chainmetric-sensorsys/model"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// Operate registers supported sensors and starts asynchronously work on reading requests.
func (d *Device) Operate() {
	if !d.active {
		return
	}

	d.reader.RegisterSensors(d.SupportedSensors()...)

	for _, request := range d.requests.Get() {
		d.actOnRequest(request)
	}

	d.reader.Process()
}

func (d *Device) actOnRequest(request *readingsRequest) {
	var (
		handler = func(readings model.SensorsReadingResults) {
			d.postReadings(request.AssetID, readings)
		}
	)

	if request.Period.Seconds() == 0 {
		d.reader.SendRequest(handler, request.Metrics...)
		delete(d.requests.data, request.ID)
		return
	}

	request.Cancel = d.reader.SubscribeReceiver(handler, request.Period, request.Metrics...)
}

func (d *Device) postReadings(assetID string, readings model.SensorsReadingResults) {
	var (
		record = models.MetricReadings{
			AssetID: assetID,
			DeviceID: d.model.ID,
			Timestamp: time.Now(),
			Values: readings,
		}
	)

	if len(readings) == 0 {
		shared.Logger.Warningf("No metrics was read for asset %s, posting is skipped", assetID)
		return
	}

	if err := blockchain.Contracts.Readings.Post(record); err != nil {
		if detectNetworkAbsence(err) {
			d.handleNetworkDisconnection(record)
		} else {
			shared.Logger.Error(errors.Wrap(err, "failed to post readings"))
		}
		return
	}

	shared.Logger.Debugf("Readings for asset %s was posted with => %s", assetID, utils.Prettify(readings))
}
