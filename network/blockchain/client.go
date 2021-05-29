package blockchain

import (
	"fmt"
	"io/ioutil"

	fabconfig "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"

	"github.com/pkg/errors"

	cnf "github.com/timoth-y/chainmetric-sensorsys/model/config"
)

// Client defines an interface for communicating with blockchain network.
type Client struct {
	wallet  *gateway.Wallet
	gateway *gateway.Gateway
	network *gateway.Network
}

var (
	client = &Client{}
	config cnf.BlockchainConfig

	// Contracts exposes blockchain network SmartContracts pool.
	Contracts = struct {
		Devices      *DevicesContract
		Assets       *AssetsContract
		Requirements *RequirementsContract
		Readings     *ReadingsContract
	} {
		Devices: &DevicesContract{},
		Assets: &AssetsContract{},
		Requirements: &RequirementsContract{},
		Readings: &ReadingsContract{},
	}
)

// Init performs initialization sequence of the blockchain client with given config.
func Init(cnf cnf.BlockchainConfig) error {
	config = cnf

	var (
		err error
		configProvider = fabconfig.FromFile(config.ConnectionConfig)
		identity *gateway.X509Identity
	)

	if client.wallet, err = gateway.NewFileSystemWallet(config.WalletPath); err != nil {
		return errors.Wrapf(err, "failed to create new wallet on %s", config.WalletPath)
	}

	if identity, err = newX509Identity(config.Identity); err != nil {
		return errors.Wrap(err, "failed to build X509 identity")
	}

	if err = client.wallet.Put(config.Identity.UserID, identity); err != nil {
		return errors.Wrap(err, "failed to put identity to wallet")
	}

	if client.gateway, err = gateway.Connect(
		gateway.WithConfig(configProvider),
		gateway.WithIdentity(client.wallet, config.Identity.UserID),
	); err != nil {
		return errors.Wrap(err, "failed to connect to blockchain gateway")
	}

	if client.network, err = client.gateway.GetNetwork(config.ChannelID); err != nil {
		return errors.Wrapf(err, "failed to create new client of channel %s", config.ChannelID)
	}

	Contracts.Assets.init()
	Contracts.Devices.init()
	Contracts.Requirements.init()
	Contracts.Readings.init()

	return nil
}

// Close closes connection to blockchain network and clears allocated resources.
func Close() {
	client.gateway.Close()
}

func newX509Identity(identity cnf.BlockchainIdentityConfig) (*gateway.X509Identity, error) {
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

