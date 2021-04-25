package blockchain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/timoth-y/iot-blockchain-contracts/models"
)

type ReadingsContract struct {
	client   *Client
	contract *gateway.Contract
}

func NewReadingsContract(client *Client) *ReadingsContract {
	return &ReadingsContract{
		client: client,
		contract: client.network.GetContract("readings"),
	}
}

func (cc *ReadingsContract) Post(readings models.MetricReadings) error {
	_, err := cc.contract.SubmitTransaction("Post", string(readings.Encode()))
	return err
}
