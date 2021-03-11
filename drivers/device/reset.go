package device

import (
	"os"

	"github.com/pkg/errors"

	"github.com/timoth-y/iot-blockchain-sensorsys/gateway/blockchain"
	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

func Reset(dc *blockchain.DevicesContract) error {
	id, is := isRegistered(); if !is {
		return nil
	}

	if err := dc.Remove(id); err != nil {
		return err
	}

	if err := os.Remove(deviceIdentityFile); err != nil {
		return errors.Wrap(err, "failed to remove device's identity file")
	}

	shared.Logger.Info("Device is been reset.")

	return nil
}
