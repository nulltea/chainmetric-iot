package device

import (
	"time"

	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/models/requests"
	"github.com/timoth-y/chainmetric-core/utils"

	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

func (d *Device) ListenRemoteCommands() error {
	var (
		contract = d.client.Contracts.Devices
	)

	if !d.active {
		return nil
	}

	if d.model == nil || len(d.model.ID) == 0 {
		return errors.New("cannot listen to commands before device registration")
	}

	go func() {
		if err := contract.ListenCommands(d.ctx, d.model.ID, func(id string, cmd models.DeviceCommand, args ...interface{}) error {
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
			shared.Logger.Error(err)
		}
	}()

	return nil
}

func (d *Device) handleBluetoothPairingCmd(cmdID string) {
	var (
		contract = d.client.Contracts.Devices
		results = requests.DeviceCommandResultsSubmitRequest{
			Status: models.DeviceCmdCompleted,
		}
	)

	if err := d.localnet.Pair(); err != nil {
		results.Status = models.DeviceCmdFailed
		results.Error = utils.StringPointer(err.Error())
		shared.Logger.Error(err)
	}

	results.Timestamp = time.Now().UTC()

	if err := contract.SubmitCommandResults(cmdID, results); err != nil {
		shared.Logger.Error(err)
	}
}
