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

	"github.com/timoth-y/chainmetric-iot/shared"
)

// AssetsContract defines interface for communication with assets-managing Smart Contract.
type AssetsContract struct {
	contract *gateway.Contract
}

// init performs initialization of the AssetsContract instance.
func (ac *AssetsContract) init() {
	ac.contract = client.network.GetContract("assets")
}

// Receive retrieves models.Asset records from blockchain ledger by a given `query`.
func (ac *AssetsContract) Receive(query requests.AssetsQuery) ([]*models.Asset, error) {
	var assets []*models.Asset

	data, err := ac.contract.EvaluateTransaction("QueryRaw", string(query.Encode())); if err != nil {
		return nil, err
	}

	if data != nil {
		if err = json.Unmarshal(data, &assets); err != nil {
			return nil, err
		}
	}

	return assets, nil
}

// Subscribe starts listening to blockchain events related to assets and triggers `action` on each event occurrence.
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
				shared.Logger.Debug("Assets blockchain event listener ended")
				return nil
			}
		}
	}
}
