package device

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/periphery"
	"github.com/timoth-y/iot-blockchain-sensorsys/drivers/sensors"
	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

func (d *Device) initHotswap() {
	var (
		ctx = context.Background()
		startTime  time.Time
	)

	ctx, d.cancelHotswap = context.WithCancel(ctx)

	go func() {
	LOOP: for {
			startTime = time.Now()

			if err := d.handleHotswap(); err != nil {
				shared.Logger.Error(errors.Wrap(err, "failed to handle hotswap"))
			}

			select {
			case <-time.After(time.Second - time.Since(startTime)):
			case <- ctx.Done():
				break LOOP
			}
		}
	}()
}

func (d *Device) handleHotswap() error {
	d.detectedI2Cs = periphery.DetectI2C(sensors.I2CAddressesRange())

	// TODO: implement hot swap changes for reader

	return nil
}
