package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/models/requests"

	"github.com/timoth-y/chainmetric-sensorsys/model"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

type DevicesContract struct {
	client   *Client
	contract *gateway.Contract
}

func NewDevicesContract(client *Client) *DevicesContract {
	return &DevicesContract{
		client: client,
	}
}

func (dc *DevicesContract) Init() {
	dc.contract = dc.client.network.GetContract("devices")
}

func (dc *DevicesContract) Retrieve(id string) (*models.Device, error) {
	resp, err := dc.contract.EvaluateTransaction("Retrieve", id); if err != nil {
		return nil, err
	}

	return models.Device{}.Decode(resp)
}

func (dc *DevicesContract) Exists(id string) (bool, error) {
	resp, err := dc.contract.EvaluateTransaction("Exists", id); if err != nil {
		return false, err
	}

	return strconv.ParseBool(string(resp))
}

func (dc *DevicesContract) UpdateSpecs(id string, specs *model.DeviceSpecs) error {
	data, err := json.Marshal(specs); if err != nil {
		return err
	}

	if  _, err = dc.contract.SubmitTransaction("Update", id, string(data)); err != nil {
		return err
	}

	return nil
}

func (dc *DevicesContract) UpdateState(id string, state models.DeviceState) error {
	data, err := json.Marshal(requests.DeviceUpdateRequest{State: &state}); if err != nil {
		return err
	}

	if  _, err = dc.contract.SubmitTransaction("Update", id, string(data)); err != nil {
		return err
	}

	return nil
}

func (dc *DevicesContract) Unbind(id string) error {
	if  _, err := dc.contract.SubmitTransaction("Unbind", id); err != nil {
		return err
	}

	return nil
}

func (dc *DevicesContract) Subscribe(ctx context.Context, event string, action func(*models.Device, string) error) error {
	reg, notifier, err := dc.contract.RegisterEvent(eventFilter("devices", event)); if err != nil {
		return err
	}

	defer dc.contract.Unregister(reg)

	for {
		select {
		case event := <-notifier:
			dev, err := models.Device{}.Decode(event.Payload); if err != nil {
				shared.Logger.Errorf("failed parse device from event payload: %v", err)
				continue
			}

			go func(d *models.Device, e string) {
				if err := action(d, e); err != nil {
					shared.Logger.Error(err)
				}
			}(dev, strings.Replace(event.EventName, "devices.", "", 1))
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
