package blockchain

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type RequirementsContract struct {
	client   *Client
	contract *gateway.Contract
}

func NewRequirementsContract(client *Client) *ReadingsContract {
	return &ReadingsContract{
		client: client,
		contract: client.network.GetContract("requirements"),
	}
}

func (cc *ReadingsContract) Subscribe(handler func()) error {
	reg, notifier, err := cc.contract.RegisterEvent("requirements.inserted")
	defer cc.contract.Unregister(reg)


	select {
	case event := <-notifier:
		fmt.Printf("Received CC event: %#v\n", string(event.Payload))
	}
	return err
}
