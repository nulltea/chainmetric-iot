package periphery

import (
	"context"
	"sync"

	"github.com/spf13/viper"
	"periph.io/x/periph/conn/i2c/i2creg"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensor"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensors"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// I2CDetectResults stores I2C identified I2C-based peripheral devices.
type I2CDetectResults map[int][]sensor.Sensor

func ScanI2C(addrs []uint16, detector func(addr uint16, bus int) (sensor.Factory, bool)) I2CDetectResults {
	var (
		detected = make(map[int][]sensor.Sensor)
		wg       = sync.WaitGroup{}
	)

	if viper.GetBool("mocks.debug_env") {
		detected[1] = []sensor.Sensor{sensors.NewI2CSensorMock(sensors.MOCK_ADDRESS, 1)}
	}

	for _, ref := range i2creg.All() {
		ctx, _ := context.WithTimeout(context.Background(), viper.GetDuration("device.i2c_scan_timeout"))
		wg.Add(1)

		go func(ctx context.Context, ref *i2creg.Ref) {
			defer wg.Done()

			bus, err := ref.Open(); if err != nil {
				shared.Logger.Error(err)
				return
			}
			defer bus.Close()

			detected[ref.Number] = make([]sensor.Sensor, 0)

			for _, addr := range addrs {
				if err := bus.Tx(addr, []byte{}, []byte{0x0}); err != nil {
					continue
				}

				if sf, ok := detector(addr, ref.Number); ok {
					detected[ref.Number] = append(detected[ref.Number], sf.Build(ref.Number))
				}

				select {
				case <- ctx.Done():
					return
				default:
					continue
				}
			}
		}(ctx, ref)
	}

	wg.Wait()

	return detected
}
