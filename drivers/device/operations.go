package device

import (
	"github.com/timoth-y/iot-blockchain-sensorsys/model"
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
		return
	}

	request.cancel = d.reader.SubscribeReceiver(handler, request.period, request.metrics...)
}

func (d *Device) postReadings(assetID string, readings model.MetricReadings) {
	// TODO: posting to blockchain
}
