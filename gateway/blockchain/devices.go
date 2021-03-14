package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/model"
	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

type DevicesContract struct {
	client   *Client
	contract *gateway.Contract
}

func NewDevicesContract(client *Client) *DevicesContract {
	return &DevicesContract{
		client: client,
		contract: client.network.GetContract("devices"),
	}
}

func (cc *DevicesContract) Exists(id string) (bool, error) {
	resp, err := cc.contract.EvaluateTransaction("Exists", id)

	if err != nil {
		return false, err
	}

	return strconv.ParseBool(string(resp))
}

func (cc *DevicesContract) UpdateSpecs(id string, specs *model.DeviceSpecs) error {
	data, err := json.Marshal(specs); if err != nil {
		return err
	}

	if  _, err = cc.contract.EvaluateTransaction("Update", id, string(data)); err != nil {
		return err
	}

	return nil
}

func (cc *DevicesContract) Remove(id string) error {
	if  _, err := cc.contract.EvaluateTransaction("Remove", id); err != nil {
		return err
	}

	return nil
}

func (cc *DevicesContract) Subscribe(ctx context.Context, event string, action func(device models.Device) error) error {
	reg, notifier, err := cc.contract.RegisterEvent(fmt.Sprintf("devices.%s", event)); if err != nil {
		return err
	}

	defer cc.contract.Unregister(reg)

	for {
		select {
		case event := <-notifier:
			dev := models.Device{}

			if err := json.Unmarshal(event.Payload, &dev); err != nil {
				shared.Logger.Errorf("failed decentralise device from event payload: %v", err)
				continue
			}

			go func(d models.Device) {
				if err := action(d); err != nil {
					shared.Logger.Error(err)
				}
			}(dev)
		case <- ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				return fmt.Errorf("timeout waiting for event devices.%s", event)
			default:
				return nil
			}
		}
	}

	return nil
}
