package blockchain

import (
	"context"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

type RequirementsContract struct {
	client   *Client
	contract *gateway.Contract
}

func NewRequirementsContract(client *Client) *RequirementsContract {
	return &RequirementsContract{
		client: client,
		contract: client.network.GetContract("requirements"),
	}
}


func (cc *RequirementsContract) Subscribe(ctx context.Context, event string, action func(*models.Requirements, string) error) error {
	reg, notifier, err := cc.contract.RegisterEvent(eventFilter("requirements", event)); if err != nil {
		return err
	}

	defer cc.contract.Unregister(reg)

	for {
		select {
		case event := <-notifier:
			req, err := models.Requirements{}.Decode(event.Payload); if err != nil {
			shared.Logger.Errorf("failed parse device from event payload: %v", err)
			continue
		}

			go func(r *models.Requirements, e string) {
				if err := action(r, e); err != nil {
					shared.Logger.Error(err)
				}
			}(req, strings.Replace(event.EventName, "requirements.", "", 1))
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
