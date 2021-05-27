package blockchain

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/timoth-y/chainmetric-core/models"
)

// ReadingsContract defines access to blockchain Smart Contract for managing metric readings.
type ReadingsContract struct {
	contract *gateway.Contract
}

// init performs initialization of the ReadingsContract instance.
func (rc *ReadingsContract) init() {
	rc.contract = client.network.GetContract("readings")
}

// Post sends models.MetricReadings record to blockchain network for processing.
func (rc *ReadingsContract) Post(readings models.MetricReadings) error {
	_, err := rc.contract.SubmitTransaction("Post", string(readings.Encode()))
	return err
}
