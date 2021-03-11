package shared

import (
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)


func InitPeriphery() {
	if _, err := host.Init(); err != nil {
		Logger.Fatal(err)
	}
}

func ScanI2CAddrs(start, end uint8) map[int][]uint8 {
	var (
		addrMap = make(map[int][]uint8)
	)

	InitPeriphery();

	for _, ref := range i2creg.All() {
		bus, err := ref.Open(); if err != nil {
			continue
		}
		addrMap[ref.Number] = make([]uint8, 0)
		for addr := start; addr <= end; addr++ {
			if err := bus.Tx(uint16(addr), []byte{}, []byte{0x0}); err == nil {
				addrMap[ref.Number] = append(addrMap[ref.Number], addr)
			}
		}
		bus.Close()
	}

	return addrMap
}


