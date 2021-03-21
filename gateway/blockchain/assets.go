package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/timoth-y/iot-blockchain-contracts/models"
	"github.com/timoth-y/iot-blockchain-contracts/models/request"

	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

type AssetsContract struct {
	client   *Client
	contract *gateway.Contract
}

func NewAssetsContract(client *Client) *AssetsContract {
	return &AssetsContract{
		client: client,
		contract: client.network.GetContract("assets"),
	}
}

func (ac *AssetsContract) Receive(query request.AssetsQuery) ([]*models.Asset, error) {
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
	reg, notifier, err := ac.contract.RegisterEvent(eventFilter("assets", event)); if err != nil {
		return err
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
				return nil
			}
		}
	}

	return nil
}
