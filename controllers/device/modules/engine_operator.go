package modules

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/utils"
	"github.com/timoth-y/chainmetric-iot/controllers/device"
	"github.com/timoth-y/chainmetric-iot/controllers/engine"
	"github.com/timoth-y/chainmetric-iot/model"
	"github.com/timoth-y/chainmetric-iot/model/events"
	"github.com/timoth-y/chainmetric-iot/network/blockchain"
	"github.com/timoth-y/chainmetric-iot/shared"
	"github.com/timoth-y/go-eventdriver"
)

// EngineOperator implements device.Module for engine.SensorsReader operating.
type EngineOperator struct {
	moduleBase
	engine *engine.SensorsReader
}

// WithEngineOperator can be used to setup EngineOperator logical device.Module onto the device.Device.
func WithEngineOperator() device.Module {
	return &EngineOperator{
		moduleBase: withModuleBase("ENGINE_OPERATOR"),
		engine: engine.NewSensorsReader(),
	}
}

func (m *EngineOperator) Start(ctx context.Context) {
	go m.Do(func() {
		if !m.trySyncWithDeviceLifecycle(ctx, m.Start) {
			return
		}

		// Listen and act on newly submitted or changed requirements:
		eventdriver.SubscribeHandler(events.RequirementsChanged, func(_ context.Context, v interface{}) error {
			if !m.engine.Active() {
				return nil
			}  // No need to act on requests before engine isn't started

			if payload, ok := v.(events.RequirementsChangedPayload); ok {
				for i := range payload.Requests {
					m.actOnRequest(ctx, payload.Requests[i])
				}

				return nil
			}

			return eventdriver.ErrIncorrectPayload
		})

		// Listen and changes in device's sensors register:
		eventdriver.SubscribeHandler(events.SensorsRegisterChanged, func(_ context.Context, v interface{}) error {
			if payload, ok := v.(events.SensorsRegisterChangedPayload); ok {
				m.engine.RegisterSensors(payload.Added...)
				m.engine.UnregisterSensors(payload.Removed...)

				// If engine wasn't started yet it is because there weren't any available sensors before.
				// If there is ones now, engine could start processing requests.
				if !m.engine.Active() && m.RegisteredSensors().NotEmpty() {
					m.engine.Run(ctx)
					m.actOnCachedRequests(ctx)
				}

				return nil
			}

			return eventdriver.ErrIncorrectPayload
		})

		// Listen and changes in parameters cache:
		eventdriver.SubscribeHandler(events.CacheChanged, func(_ context.Context, _ interface{}) error {
			m.actOnCachedRequests(ctx)
			return nil
		})

		if m.waitUntilSensorsDetected() {
			m.engine.RegisterSensors(m.RegisteredSensors().ToList()...)
			m.engine.Run(ctx)
		}
	})
}


func (m *EngineOperator) actOnRequest(ctx context.Context, request *model.SensorsReadingRequest) {
	if request.IsProcessed() {
		return
	}

	var (
		handler = func(readings engine.ReadingResults) {
			m.postReadings(request.AssetID, readings)
			eventdriver.EmitEvent(ctx, events.RequestHandled, nil)
		}
	)

	// Handle one-time request
	if request.Period.Seconds() == 0 {
		m.engine.SendRequest(handler, request.Metrics...)
		m.RemoveRequirementsFromCache(request.ID)
		return
	}

	// Otherwise subscribe receiver with given period of readings:
	request.SetCancel(m.engine.SubscribeReceiver(ctx, handler, request.Period, request.Metrics...))
}

func (m *EngineOperator) actOnCachedRequests(ctx context.Context) {
	for _, request := range m.GetCachedRequirements() {
		m.actOnRequest(ctx, request)
	}
}

func (m *EngineOperator) postReadings(assetID string, readings engine.ReadingResults) {
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
			shared.Logger.Error(errors.Wrapf(err, "failed to post readings with id %s%s", utils.Hash(record.AssetID),
				utils.Hash(string(record.Encode()))))
		}
		return
	}

	shared.Logger.Debugf("Readings for asset %s was posted with => %s", assetID, utils.Prettify(readings))
}
