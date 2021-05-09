package device

import (
	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-core/models"
)

func (d *Device) handleLocationTracking() {
	var (
		contract = d.client.Contracts.Devices
	)

	d.localnet.Channels.Geo.Subscribe(d.ctx, func(location models.Location) error {
		return errors.Wrap(contract.UpdateLocation(d.model.ID, location), "failed to update device location")
	})
}
