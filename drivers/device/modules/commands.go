package modules

import (
	"context"
	"sync"
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

// RemoteCommandsHandler defines device.Device module for remote commands handling.
type RemoteCommandsHandler struct {
	dev  *dev.Device
	once *sync.Once
}

// WithRemoteCommandsHandler can be used to setup RemoteCommandsHandler module for the device.Device.
func WithRemoteCommandsHandler() Module {
	return &RemoteCommandsHandler{
		once: &sync.Once{},
	}
}

func (m *RemoteCommandsHandler) Setup(device *dev.Device) error {
	m.dev = device

	return nil
}

func (m *RemoteCommandsHandler) Start(ctx context.Context) {
	go func() {
		if err := blockchain.Contracts.Devices.ListenCommands(ctx,
			m.dev.ID(), func(id string, cmd models.DeviceCommand, args ...interface{}) error {
				switch cmd {
				case models.DevicePauseCmd:
				case models.DeviceResumeCmd:
				case models.DevicePairingCmd:
					m.handleBluetoothPairingCmd(id)
				default:
					shared.Logger.Error(errors.Errorf("command '%s' is not supported", cmd))
				}
				return nil
			}); err != nil {
			shared.Logger.Error(errors.Wrap(err, "failed to subscribe to device remote commands"))
		}
	}()
}

func (d *RemoteCommandsHandler) handleBluetoothPairingCmd(cmdID string) {
	var (
		results = requests.DeviceCommandResultsSubmitRequest{
			Status: models.DeviceCmdCompleted,
		}
	)

	if err := localnet.Pair(); err != nil {
		results.Status = models.DeviceCmdFailed
		results.Error = utils.StringPointer(err.Error())
		shared.Logger.Error(err)
	}

	results.Timestamp = time.Now().UTC()

	if err := blockchain.Contracts.Devices.SubmitCommandResults(cmdID, results); err != nil {
		shared.Logger.Error(err)
	}
}
