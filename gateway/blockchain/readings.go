package blockchain

import (
	"encoding/json"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"

	"github.com/timoth-y/iot-blockchain-contracts/models"

	"github.com/timoth-y/iot-blockchain-sensorsys/model"
)

type ReadingsContract struct {
	client   *Client
	contract *gateway.Contract
}

func NewSensorReadingsContract(client *Client) *ReadingsContract {
	return &ReadingsContract{
		client: client,
		contract: client.network.GetContract("engine"),
	}
}

func (cc *ReadingsContract) Post(assetID string, values model.MetricReadingsResults) error {
	readings := models.MetricReadings{
		AssetID: assetID,
		Values: values,
	}

	data, err := json.Marshal(readings); if err != nil {
		return err
	}

	_, err = cc.contract.SubmitTransaction("Insert", string(data))

	return err
}
