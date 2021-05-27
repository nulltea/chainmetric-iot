package device

import (
	"time"

	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/models/requests"
	"github.com/timoth-y/chainmetric-core/utils"
	"github.com/timoth-y/chainmetric-sensorsys/network/blockchain"
	"github.com/timoth-y/chainmetric-sensorsys/network/localnet"

	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

func (d *Device) ListenRemoteCommands() error {
	if !d.active {
		return nil
	}

	if d.model == nil || len(d.model.ID) == 0 {
		return errors.New("cannot listen to commands before device registration")
	}

	go func() {
		if err := blockchain.Contracts.Devices.ListenCommands(
			d.ctx, d.model.ID, func(id string, cmd models.DeviceCommand, args ...interface{}) error {
			switch cmd {
			case models.DevicePauseCmd:
			case models.DeviceResumeCmd:
			case models.DevicePairingCmd:
				d.handleBluetoothPairingCmd(id)
			default:
				shared.Logger.Error(errors.Errorf("command '%s' is not supported", cmd))
			}
			return nil
		}); err != nil {
			shared.Logger.Error(errors.Wrap(err, "failed to subscribe to device remote commands"))
		}
	}()

	return nil
}

func (d *Device) handleBluetoothPairingCmd(cmdID string) {
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
