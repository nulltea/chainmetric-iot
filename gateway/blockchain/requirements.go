package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

type RequirementsContract struct {
	client   *Client
	contract *gateway.Contract
}

func NewRequirementsContract(client *Client) *RequirementsContract {
	return &RequirementsContract{
		client: client,
	}
}

func (rc *RequirementsContract) Init() {
	rc.contract = rc.client.network.GetContract("requirements")
}

func (rc *RequirementsContract) ReceiveFor(assets []string) ([]*models.Requirements, error) {
	request, _ := json.Marshal(assets)

	data, err := rc.contract.EvaluateTransaction("ForAssets", string(request)); if err != nil {
		return nil, err
	}

	var requirements []*models.Requirements; if err = json.Unmarshal(data, &requirements); err != nil {
		return nil, err
	}

	return requirements, nil
}

func (rc *RequirementsContract) Subscribe(ctx context.Context, event string, action func(*models.Requirements, string) error) error {
	reg, notifier, err := rc.contract.RegisterEvent(eventFilter("requirements", event)); if err != nil {
		return err
	}

	defer rc.contract.Unregister(reg)

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
				shared.Logger.Debug("Requirements blockchain event listener ended.")
				return nil
			}
		}
	}

	return nil
}
