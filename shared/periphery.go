package shared

import (
	"context"
	"sync"
	"time"

	"github.com/timoth-y/iot-blockchain-contracts/shared"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

func InitPeriphery() {
	if _, err := host.Init(); err != nil {
		Logger.Fatal(err)
	}
}

func ScanI2CAddrs(start, end uint16) map[int][]uint16 {
	var (
		addrMap = make(map[int][]uint16)
		wg = sync.WaitGroup{}
	)

	InitPeriphery()

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


