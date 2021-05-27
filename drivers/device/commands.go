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

	if d.state == nil || len(d.ID()) == 0 {
		return errors.New("cannot listen to commands before device registration")
	}

	go func() {
		if err := blockchain.Contracts.Devices.ListenCommands(
			d.ctx, d.state.ID, func(id string, cmd models.DeviceCommand, args ...interface{}) error {
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


