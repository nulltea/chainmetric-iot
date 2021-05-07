package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/models/requests"

	"github.com/timoth-y/chainmetric-sensorsys/model"
	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// DevicesContract defines access to blockchain Smart Contract for managing device.
type DevicesContract struct {
	client   *Client
	contract *gateway.Contract
}

// NewDevicesContract constructs new DevicesContract instance.
func NewDevicesContract(client *Client) *DevicesContract {
	return &DevicesContract{
		client: client,
	}
}

// Init performs initialization of DevicesContract.
func (dc *DevicesContract) Init() {
	dc.contract = dc.client.network.GetContract("devices")
}

// Retrieve fetches models.Device from the blockchain ledger.
func (dc *DevicesContract) Retrieve(id string) (*models.Device, error) {
	resp, err := dc.contract.EvaluateTransaction("Retrieve", id)
	if err != nil {
		return nil, err
	}

	return models.Device{}.Decode(resp)
}

// Exists verify whether the device with `id` exists on the blockchain ledger.
func (dc *DevicesContract) Exists(id string) (bool, error) {
	resp, err := dc.contract.EvaluateTransaction("Exists", id)
	if err != nil {
		return false, err
	}

	return strconv.ParseBool(string(resp))
}

// UpdateSpecs updates device specification on the blockchain ledger.
func (dc *DevicesContract) UpdateSpecs(id string, specs *model.DeviceSpecs) error {
	data, err := json.Marshal(specs)
	if err != nil {
		return err
	}

	if _, err = dc.contract.SubmitTransaction("Update", id, string(data)); err != nil {
		return err
	}

	return nil
}

// UpdateSpecs updates device state on the blockchain ledger.
func (dc *DevicesContract) UpdateState(id string, state models.DeviceState) error {
	data, err := json.Marshal(requests.DeviceUpdateRequest{State: &state})
	if err != nil {
		return err
	}

	if _, err = dc.contract.SubmitTransaction("Update", id, string(data)); err != nil {
		return err
	}

	return nil
}

// Unbind removes device from the blockchain ledger.
func (dc *DevicesContract) Unbind(id string) error {
	if _, err := dc.contract.SubmitTransaction("Unbind", id); err != nil {
		return err
	}

	return nil
}

// Subscribe subscribes and starts listening for device related events on the blockchain network.
func (dc *DevicesContract) Subscribe(
	ctx context.Context, event string,
	action func(*models.Device, string) error,
) error {
	var (
		eventFilter = eventFilter("devices", event)
		reg, notifier, err = dc.contract.RegisterEvent(eventFilter)
	)

	if err != nil {
		return errors.Wrapf(err, "failed to subscribe to '%s' events", eventFilter)
	}

	defer dc.contract.Unregister(reg)

	for {
		select {
		case event := <-notifier:
			dev, err := models.Device{}.Decode(event.Payload)
			if err != nil {
				shared.Logger.Errorf("failed parse device from event payload: %v", err)
				continue
			}

			go func(d *models.Device, e string) {
				if err := action(d, e); err != nil {
					shared.Logger.Error(err)
				}
			}(dev, strings.Replace(event.EventName, "devices.", "", 1))
		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				return fmt.Errorf("timeout waiting for event devices.%s", event)
			default:
				shared.Logger.Debug("Device blockchain event listener ended.")
				return nil
			}
		}
	}

	return nil
}

// ListenCommands subscribes and starts listening for device commands from the blockchain network.
func (dc *DevicesContract) ListenCommands(
	ctx context.Context, deviceID string,
	handler func(id string, cmd models.DeviceCommand, args ...interface{}) error,
) error {
	var (
		eventKey = fmt.Sprintf("devices.%s.command", deviceID)
		reg, notifier, err = dc.contract.RegisterEvent(eventKey)
	)

	if err != nil {
		return errors.Wrapf(err, "failed to subscribe to '%s' events", eventKey)
	}

	defer dc.contract.Unregister(reg)

	for {
		select {
		case event := <-notifier:
			cmd, err := requests.DeviceCommandEventPayload{}.Decode(event.Payload)
			if err != nil {
				shared.Logger.Errorf("failed parse device from event payload: %v", err)
				continue
			}

			go func(c *requests.DeviceCommandEventPayload) {
				if err := handler(c.ID, c.Command, c.Args...); err != nil {
					shared.Logger.Error(err)
				}
			}(cmd)
		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				return errors.Errorf("timeout waiting for event %s", eventKey)
			default:
				shared.Logger.Debug("Device blockchain event listener ended.")
				return nil
			}
		}
	}

	return nil
}

// SubmitCommandResults submits command execution results to log them in the blockchain ledger.
func (dc *DevicesContract) SubmitCommandResults(id string, req requests.DeviceCommandResultsSubmitRequest) error {
	if _, err := dc.contract.SubmitTransaction("SubmitCommandResults", id, string(req.Encode())); err != nil {
		return errors.Wrapf(err, "failed to submit command execution results for id '%s'", id)
	}

	return nil
}
