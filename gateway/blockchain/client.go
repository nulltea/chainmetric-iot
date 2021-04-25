package blockchain

import (
	"fmt"
	"io/ioutil"

	fabconfig "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"

	"github.com/pkg/errors"

	"github.com/timoth-y/chainmetric-sensorsys/model/config"
)

type Client struct {
	wallet  *gateway.Wallet
	gateway *gateway.Gateway
	network *gateway.Network

	Contracts contracts
}

type contracts struct {
	Devices      *DevicesContract
	Assets       *AssetsContract
	Requirements *RequirementsContract
	Readings     *ReadingsContract
}

func NewBlockchainClient() *Client {
	return &Client{}
}

func (bc *Client) Init(config config.BlockchainConfig) (err error) {
	configProvider := fabconfig.FromFile(config.ConnectionConfig)

	bc.wallet, err = gateway.NewFileSystemWallet(config.WalletPath)
	if err != nil {
		err = errors.Wrapf(err, "failed to create new wallet on %s", config.WalletPath)
		return
	}

	identity, err := newX509Identity(config.Identity); if err != nil {
		err = errors.Wrap(err, "failed to build X509 identity")
		return
	}

	if err = bc.wallet.Put(config.Identity.UserID, identity); err != nil {
		err = errors.Wrap(err, "failed to put identity to wallet")
		return
	}

	bc.gateway, err = gateway.Connect(
		gateway.WithConfig(configProvider),
		gateway.WithIdentity(bc.wallet, config.Identity.UserID),
	); if err != nil {
		return errors.Wrap(err, "failed to connect to blockchain gateway")
	}

	bc.network, err = bc.gateway.GetNetwork(config.ChannelID); if err != nil {
		err = errors.Wrapf(err, "failed to create new client of channel %s", config.ChannelID)
		return
	}

	bc.Contracts = contracts{
		Devices:      NewDevicesContract(bc),
		Assets:       NewAssetsContract(bc),
		Requirements: NewRequirementsContract(bc),
		Readings:     NewReadingsContract(bc),
	}

	return
}

func (bc *Client) Close() {
	bc.gateway.Close()
}

func newX509Identity(identity config.BlockchainIdentityConfig) (*gateway.X509Identity, error) {
	cert, err := ioutil.ReadFile(identity.Certificate); if err != nil {
		return nil, err
	}

	key, err := ioutil.ReadFile(identity.PrivateKey); if err != nil {
		return nil, err
	}

	return gateway.NewX509Identity(identity.MspID, string(cert), string(key)), nil
}

func eventFilter(prefix, action string) string {
	if len(prefix) == 0 {
		return action
	}

	if action == "*" {
		action = "[a-zA-Z]+"
	}

	return fmt.Sprintf(`%s\.%s`, prefix, action)
}

