package device

import (
	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

func (d *Device) handleLocationTracking() {
	var (
		contract = d.client.Contracts.Devices
	)

	d.localnet.Channels.Geo.Subscribe(d.ctx, func(location models.Location) error {
		if err := errors.Wrap(contract.UpdateLocation(d.model.ID, location), "failed to update device location"); err != nil {
			return err
		}

		shared.Logger.Debugf("Device location was updated via Bluetooth tethering: %s", location.Name)
		return nil
	})
}
