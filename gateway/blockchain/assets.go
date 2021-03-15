package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/timoth-y/iot-blockchain-contracts/models"

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

func (cc *AssetsContract) Receive() ([]*models.Asset, error) {
	data, err := cc.contract.EvaluateTransaction("List"); if err != {
		return nil, err
	}

	assets := []*models.Asset{}

	if err = json.Unmarshal(data, assets); err != nil {
		return nil, err
	}

	return assets, nil
}

func (cc *AssetsContract) Subscribe(ctx context.Context, event string, action func(*models.Asset, string) error) error {
	reg, notifier, err := cc.contract.RegisterEvent(eventFilter("assets", event)); if err != nil {
		return err
	}

	defer cc.contract.Unregister(reg)

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
