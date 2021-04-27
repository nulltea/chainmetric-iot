package periphery

import (
	"context"
	"sync"
	"time"

	"github.com/spf13/viper"
	"periph.io/x/periph/conn/i2c/i2creg"

	"github.com/timoth-y/chainmetric-sensorsys/drivers/sensors"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// I2CDetectResults stores I2C identified I2C-based peripheral devices.
type I2CDetectResults map[int][]uint16

// I2CDetectKey defines constants needed to identify I2C-based peripheral devices.
type I2CDetectKey struct {
	Address    uint16
	IDRegister uint16
	ID         uint16
}

func NewI2CDetectKey(addr, reg, id uint16) I2CDetectKey {
	return I2CDetectKey{
		Address: addr,
		IDRegister: reg,
		ID: reg,
	}
}

func DetectI2C(start, end uint16) I2CDetectResults {
	var (
		addrMap = make(map[int][]uint16)
		wg = sync.WaitGroup{}
	)

	if viper.GetBool("mocks.debug_env") {
		addrMap[1] = []uint16{sensors.MOCK_ADDRESS}
	}

	for _, ref := range i2creg.All() {
		ctx, _ := context.WithTimeout(context.Background(), 1 * time.Second)
		wg.Add(1)

		go func(ctx context.Context, ref *i2creg.Ref) {
			defer wg.Done()

			bus, err := ref.Open(); if err != nil {
				shared.Logger.Error(err)
				return
			}
			defer bus.Close()

			addrMap[ref.Number] = make([]uint16, 0)

			for addr := start; addr <= end; addr++ {
				if err := bus.Tx(addr, []byte{}, []byte{0x0}); err == nil {
					addrMap[ref.Number] = append(addrMap[ref.Number], addr)
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

	return addrMap
}
