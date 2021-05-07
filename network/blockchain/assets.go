package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/models/requests"

	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

type AssetsContract struct {
	client   *Client
	contract *gateway.Contract
}

func NewAssetsContract(client *Client) *AssetsContract {
	return &AssetsContract{
		client: client,
	}
}

func (ac *AssetsContract) Init() {
	ac.contract = ac.client.network.GetContract("assets")
}

func (ac *AssetsContract) Receive(query requests.AssetsQuery) ([]*models.Asset, error) {
	data, err := ac.contract.EvaluateTransaction("QueryRaw", string(query.Encode())); if err != nil {
		return nil, err
	}

	var assets []*models.Asset

	if err = json.Unmarshal(data, &assets); err != nil {
		return nil, err
	}

	return assets, nil
}

func (ac *AssetsContract) Subscribe(ctx context.Context, event string, action func(*models.Asset, string) error) error {
	var (
		eventFilter = eventFilter("assets", event)
		reg, notifier, err = ac.contract.RegisterEvent(eventFilter)
	)

	if err != nil {
		return errors.Wrapf(err, "failed to subscribe to '%s' events", eventFilter)
	}

	defer ac.contract.Unregister(reg)

	for {
		select {
		case event := <-notifier:
			asset, err := models.Asset{}.Decode(event.Payload); if err != nil {
				shared.Logger.Errorf("failed parse asset from event payload: %v", err)
				continue
			}
			go func(a *models.Asset, e string) {
				if err := action(a, e); err != nil {
					shared.Logger.Error(err)
				}
			}(asset, strings.Replace(event.EventName, "assets.", "", 1))
		case <- ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				return fmt.Errorf("timeout waiting for event assets.%s", event)
			default:
				shared.Logger.Debug("Assets blockchain event listener ended.")
				return nil
			}
		}
	}

	return nil
}
