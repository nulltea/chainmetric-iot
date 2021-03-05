package gateway

import (
	"io/ioutil"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	fabconfig "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"

	"github.com/pkg/errors"

	"sensorsys/config"
)

type BlockchainClient struct {
	wallet  *gateway.Wallet
	gateway *gateway.Gateway
	network *gateway.Network
	sdk     *fabsdk.FabricSDK
	client  *channel.Client
}

func NewBlockchainClient() *BlockchainClient {
	return &BlockchainClient {

	}
}

func (bc *BlockchainClient) Init(config config.BlockchainConfig) (err error) {
	configProvider := fabconfig.FromFile(config.ConnectionConfig)
	bc.sdk, err = fabsdk.New(configProvider)
	if err != nil {
		err = errors.Wrap(err, "failed to create new SDK")
		return
	}

	bc.client, err = channel.New(bc.sdk.ChannelContext(config.ChannelID,
		fabsdk.WithUser(config.Identity.UserID),
		fabsdk.WithOrg(config.Identity.OrgID),
	)); if err != nil {
		err = errors.Wrapf(err, "failed to create new client of channel %s", config.ChannelID)
		return
	}

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
		return errors.Wrap(err, "InitConnectionFor: connect to gateway")
	}

	bc.network, err = bc.gateway.GetNetwork(config.ChannelID); if err != nil {
		err = errors.Wrapf(err, "failed to create new client of channel %s", config.ChannelID)
		return
	}

	return
}

func (bc *BlockchainClient) Close() {
	bc.sdk.Close()
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
