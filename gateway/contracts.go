package gateway

import (
	"encoding/json"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"

	"sensorsys/model"
)

type SensorReadingsContract struct {
	client   *BlockchainClient
	contract *gateway.Contract
}

func NewSensorReadingsContract(client *BlockchainClient) *SensorReadingsContract {
	return &SensorReadingsContract{
		client: client,
		contract: client.network.GetContract("sensor-readings"),
	}
}

func (cc *SensorReadingsContract) PostReadings(reading model.MetricReadingsResults) error {
	data, err := json.Marshal(reading); if err != nil {
		return err
	}

	_, err = cc.contract.SubmitTransaction("Post", string(data))

	return err
}
