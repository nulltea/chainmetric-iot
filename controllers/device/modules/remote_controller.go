package modules

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/models/requests"
	"github.com/timoth-y/chainmetric-core/utils"
	"github.com/timoth-y/chainmetric-iot/controllers/device"
	"github.com/timoth-y/chainmetric-iot/model/events"
	"github.com/timoth-y/chainmetric-iot/network/blockchain"
	"github.com/timoth-y/chainmetric-iot/network/localnet"
	"github.com/timoth-y/chainmetric-iot/shared"
	"github.com/timoth-y/go-eventdriver"
)

// RemoteController implements device.Module for device.Device remote commands handling.
type RemoteController struct {
	moduleBase
}

// WithRemoteController can be used to setup RemoteController logical device.Module onto the device.Device.
func WithRemoteController() device.Module {
	return &RemoteController{
		moduleBase: withModuleBase("REMOTE_CONTROLLER"),
	}
}

func (m *RemoteController) Start(ctx context.Context) {
	go m.Do(func() {
		if !m.trySyncWithDeviceLifecycle(ctx, m.Start) {
			return
		}

		if err := blockchain.Contracts.Devices.ListenCommands(ctx, m.ID(),
			func(id string, cmd models.DeviceCommand, args ...interface{}) error {
				switch cmd {
				case models.DevicePauseCmd:
				case models.DeviceResumeCmd:
				case models.DevicePairingCmd:
					m.handleBluetoothPairingCmd(ctx, id)
				default:
					shared.Logger.Error(errors.Errorf("command '%s' is not supported", cmd))
				}
				return nil
			},
		); err != nil {
			shared.Logger.Error(errors.Wrap(err, "failed to subscribe to device remote commands"))
		}
	})
}

func (m *RemoteController) handleBluetoothPairingCmd(ctx context.Context, cmdID string) {
	var (
		results = requests.DeviceCommandResultsSubmitRequest{
			Status: models.DeviceCmdCompleted,
		}
	)

	eventdriver.EmitEvent(ctx, events.BluetoothPairingStarted, nil)

	if err := localnet.Pair(ctx); err != nil && errors.Cause(err) != context.DeadlineExceeded {
		results.Status = models.DeviceCmdFailed
		results.Error = utils.StringPointer(err.Error())
		shared.Logger.Error(err)
	}

	results.Timestamp = time.Now().UTC()

	if err := blockchain.Contracts.Devices.SubmitCommandResults(cmdID, results); err != nil {
		shared.Logger.Error(err)
	}
}
