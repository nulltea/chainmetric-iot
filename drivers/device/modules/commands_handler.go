package modules

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/models/requests"
	"github.com/timoth-y/chainmetric-core/utils"
	dev "github.com/timoth-y/chainmetric-sensorsys/drivers/device"
	"github.com/timoth-y/chainmetric-sensorsys/network/blockchain"
	"github.com/timoth-y/chainmetric-sensorsys/network/localnet"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// RemoteCommandsHandler implements device.Module for device.Device remote commands handling.
type RemoteCommandsHandler struct {
	moduleBase
}

// WithRemoteCommandsHandler can be used to setup RemoteCommandsHandler logical device.Module onto the device.Device.
func WithRemoteCommandsHandler() dev.Module {
	return &RemoteCommandsHandler{
		moduleBase: withModuleBase("remote_commands_handler"),
	}
}

func (m *RemoteCommandsHandler) Start(ctx context.Context) {
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

func (m *RemoteCommandsHandler) handleBluetoothPairingCmd(ctx context.Context, cmdID string) {
	var (
		results = requests.DeviceCommandResultsSubmitRequest{
			Status: models.DeviceCmdCompleted,
		}
	)

	if err := localnet.Pair(ctx); err != nil {
		results.Status = models.DeviceCmdFailed
		results.Error = utils.StringPointer(err.Error())
		shared.Logger.Error(err)
	}

	results.Timestamp = time.Now().UTC()

	if err := blockchain.Contracts.Devices.SubmitCommandResults(cmdID, results); err != nil {
		shared.Logger.Error(err)
	}
}
