package blockchain

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type DevicesContract struct {
	client   *Client
	contract *gateway.Contract
}

func NewDevicesContract(client *Client) *DevicesContract {
	return &DevicesContract{
		client: client,
		contract: client.network.GetContract("devices"),
	}
}

func (cc *DevicesContract) Receive() {

}

func (cc *DevicesContract) Subscribe(handler func()) error {
	reg, notifier, err := cc.contract.RegisterEvent("devices.inserted")
	defer cc.contract.Unregister(reg)


	select {
	case event := <-notifier:
		fmt.Printf("Received CC event: %#v\n", string(event.Payload))
	}
	return err
}
