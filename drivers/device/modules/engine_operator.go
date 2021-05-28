package modules

import (
	"context"
	"sync"
	"time"

	fabricStatus "github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/utils"
	dev "github.com/timoth-y/chainmetric-sensorsys/drivers/device"
	"github.com/timoth-y/chainmetric-sensorsys/engine"
	"github.com/timoth-y/chainmetric-sensorsys/model"
	"github.com/timoth-y/chainmetric-sensorsys/model/events"
	"github.com/timoth-y/chainmetric-sensorsys/network/blockchain"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
	"github.com/timoth-y/go-eventdriver"
)

// EngineOperator defines device.Device module for engine.SensorsReader operating.
type EngineOperator struct {
	*dev.Device
	once   *sync.Once
	engine *engine.SensorsReader
}

// WithEngineOperator can be used to setup EngineOperator module for the device.Device.
func WithEngineOperator() Module {
	return &EngineOperator{
		once:   &sync.Once{},
		engine: engine.NewSensorsReader(),
	}
}

func (m *EngineOperator) Setup(device *dev.Device) error {
	m.Device = device

	return nil
}

func (m *EngineOperator) Start(ctx context.Context) {
	m.once.Do(func() {
		// TODO
		m.engine.RegisterSensors(d.SupportedSensors()...)

		// Act on already cached requirements.
		for _, request := range m.GetCachedRequirements() {
			m.actOnRequest(request)
		}

		// Listen and act on newly submitted or changed requirements.
		eventdriver.SubscribeHandler(events.RequirementsSubmitted, func(_ context.Context, v interface{}) error {
			if payload, ok := v.(events.RequirementsSubmittedPayload); ok {
				for i := range payload.Requests {
					m.actOnRequest(payload.Requests[i])
				}

				return nil
			}

			return eventdriver.ErrIncorrectPayload
		})

		// Finally start working on requests
		go m.engine.Process(ctx)
	})
}


func (m *EngineOperator) actOnRequest(request model.SensorsReadingRequest) {
	var (
		handler = func(readings model.SensorsReadingResults) {
			m.postReadings(request.AssetID, readings)
		}
	)

	// Handle one-time request
	if request.Period.Seconds() == 0 {
		m.engine.SendRequest(handler, request.Metrics...)
		m.RemoveRequirementsFromCache(request.ID)
		return
	}

	// Otherwise subscribe receiver with given period of readings.
	request.Cancel = m.engine.SubscribeReceiver(handler, request.Period, request.Metrics...)
}

func (m *EngineOperator) postReadings(assetID string, readings model.SensorsReadingResults) {
	var (
		ctx = context.Background()
		record = models.MetricReadings{
			AssetID:   assetID,
			DeviceID:  m.ID(),
			Timestamp: time.Now(),
			Values:    readings,
		}
	)

	if len(readings) == 0 {
		shared.Logger.Warningf("No metrics was read for asset %s, posting is skipped", assetID)
		return
	}

	if err := blockchain.Contracts.Readings.Post(record); err != nil {
		if detectNetworkAbsence(err) {
			eventdriver.EmitEvent(ctx, events.MetricReadingsPostFailed, events.MetricReadingsPostFailedPayload{
				MetricReadings: record,
				Error: err,
			})
		} else {
			shared.Logger.Error(errors.Wrap(err, "failed to post readings"))
		}
		return
	}

	shared.Logger.Debugf("Readings for asset %s was posted with => %s", assetID, utils.Prettify(readings))
}
