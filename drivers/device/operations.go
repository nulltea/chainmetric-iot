package device

import (
	"time"

	fabricStatus "github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	"github.com/pkg/errors"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/storage"
	"github.com/timoth-y/iot-blockchain-sensorsys/model"
	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

func (d *Device) Operate() {
	d.reader.RegisterSensors(d.SupportedSensors()...)

	for _, request := range d.requests.Get() {
		d.actOnRequest(request)
	}

	d.reader.Process()
}

func (d *Device) actOnRequest(request *readingsRequest) {
	var (
		handler = func(readings model.MetricReadings) {
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

func (d *Device) postReadings(assetID string, readings model.MetricReadings) {
	var (
		contract = d.client.Contracts.Readings
		reading = models.MetricReadings{
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

	if err := contract.Post(reading); err != nil {
		d.handleReadingsPostErrors(err, reading)
		return
	}

	shared.Logger.Debugf("Readings for asset %s was posted with => %s", assetID, shared.Prettify(readings))
}

func (d *Device) handleReadingsPostErrors(err error, readings models.MetricReadings) {
	if status, ok := fabricStatus.FromError(err); ok {
		switch status.Group {
		case fabricStatus.DiscoveryServerStatus:
			if err := storage.CacheReadings(readings); err != nil {
				shared.Logger.Error(errors.Wrap(err, "failed to cache readings while network connection absence"))
				return
			}

			shared.Logger.Warning(
				"Detected network connection absence, cached readings for asset %s to post later => %s",
				shared.Prettify(readings),
			)
		}
	}

	shared.Logger.Error(errors.Wrap(err, "failed to post readings"))
}
