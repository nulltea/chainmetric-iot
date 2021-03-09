package blockchain

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
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

func (cc *AssetsContract) Receive() {

}

func (cc *AssetsContract) Subscribe(handler func()) error {
	reg, notifier, err := cc.contract.RegisterEvent("requirements.inserted")
	defer cc.contract.Unregister(reg)


	select {
	case event := <-notifier:
		fmt.Printf("Received CC event: %#v\n", string(event.Payload))
	}
	return err
}
