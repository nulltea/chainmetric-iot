package utils

import (
	"log"

	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/host"
)

func ScanI2CAddrs(start, end uint8) map[int][]uint8 {
	var (
		addrMap = make(map[int][]uint8)
	)

	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

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


