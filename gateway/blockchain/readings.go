package blockchain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/timoth-y/chainmetric-core/models"
)

type ReadingsContract struct {
	client   *Client
	contract *gateway.Contract
}

func NewReadingsContract(client *Client) *ReadingsContract {
	return &ReadingsContract{
		client: client,
	}
}

func (rc *ReadingsContract) Init() {
	rc.contract = rc.client.network.GetContract("readings")
}

func (rc *ReadingsContract) Post(readings models.MetricReadings) error {
	_, err := rc.contract.SubmitTransaction("Post", string(readings.Encode()))
	return err
}
