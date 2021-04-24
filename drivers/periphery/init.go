package periphery

import (
	"periph.io/x/periph/host"

	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

func Init() {
	if _, err := host.Init(); err != nil {
		shared.Logger.Fatal(err)
	}
}

