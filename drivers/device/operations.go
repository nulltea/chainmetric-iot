package device

import (
	"time"

	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-sensorsys/model"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// Operate registers supported sensors and starts asynchronously work on reading requests.
func (d *Device) Operate() {
	d.reader.RegisterSensors(d.SupportedSensors()...)

	for _, request := range d.requests.Get() {
		d.actOnRequest(request)
	}

	d.reader.Process()
}

func (d *Device) actOnRequest(request *readingsRequest) {
	var (
		handler = func(readings model.SensorsReadingResults) {
			d.postReadings(request.assetID, readings)
		}
	)

	if request.period.Seconds() == 0 {
		d.reader.SendRequest(handler, request.metrics...)
		delete(d.requests.data, request.id)
		return
	}

	request.cancel = d.reader.SubscribeReceiver(handler, request.period, request.metrics...)
}

func (d *Device) postReadings(assetID string, readings model.SensorsReadingResults) {
	var (
		contract = d.client.Contracts.Readings
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

	if err := contract.Post(record); err != nil {
		if detectNetworkAbsence(err) {
			d.handleNetworkDisconnection(record)
		} else {
			shared.Logger.Error(errors.Wrap(err, "failed to post readings"))
		}
		return
	}

	shared.Logger.Debugf("Readings for asset %s was posted with => %s", assetID, shared.Prettify(readings))
}
