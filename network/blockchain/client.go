package blockchain

import (
	"fmt"
	"io/ioutil"

	fabconfig "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/spf13/viper"

	"github.com/pkg/errors"
)

// Client defines an interface for communicating with blockchain network.
type Client struct {
	wallet  *gateway.Wallet
	gateway *gateway.Gateway
	network *gateway.Network
}

var (
	client     = &Client{}
	connConfig *viper.Viper
	userID     string

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
func Init() (err error) {
	connConfig = viper.New()
	connConfig.SetConfigFile(viper.GetString("blockchain.connection_config"))
	if err := connConfig.ReadInConfig(); err != nil {
		return errors.Wrapf(
			err, "failed to get connection config from path '%s'",
			viper.GetString("blockchain.connection_config"),
		)
	}

	if client.wallet, err = gateway.NewFileSystemWallet(viper.GetString("blockchain.wallet_path")); err != nil {
		return errors.Wrapf(
			err, "failed to create new wallet on %s",
			viper.GetString("blockchain.wallet_path"),
		)
	}

	if !connConfig.IsSet("client.organization") {
		return errors.New("connection config missing 'client.organization' property")
	}

	identity := gateway.NewX509Identity(connConfig.GetString("client.organization"), "", "")

	if payload, err := ioutil.ReadFile(viper.GetString("blockchain.identity.certificate")); err != nil {
		return errors.Wrapf(
			err, "failed to load certificate from path: %s",
			viper.GetString("blockchain.identity.certificate"),
		)
	} else {
		identity.Credentials.Certificate = string(payload)
	}

	if payload, err := ioutil.ReadFile(viper.GetString("blockchain.identity.private_key")); err != nil {
		return errors.Wrapf(
			err, "failed to load private key from path: %s",
			viper.GetString("blockchain.identity.private_key"),
		)
	} else {
		identity.Credentials.Key = string(payload)
	}

	userID = connConfig.GetString("x-device-userID")

	if err = client.wallet.Put(userID, identity); err != nil {
		return errors.Wrap(err, "failed to put identity to wallet")
	}

	if client.gateway, err = gateway.Connect(
		gateway.WithConfig(fabconfig.FromFile(viper.GetString("blockchain.connection_config"))),
		gateway.WithIdentity(client.wallet, userID),
	); err != nil {
		return errors.Wrap(err, "failed to connect to blockchain gateway")
	}

	if !connConfig.IsSet("client.channel") {
		return errors.New("connection config missing 'client.organization' property")
	}

	if client.network, err = client.gateway.GetNetwork(connConfig.GetString("client.channel")); err != nil {
		return errors.Wrapf(
			err, "failed to create new client of channel %s",
			connConfig.GetString("client.channel"),
		)
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

func eventFilter(prefix, action string) string {
	if len(prefix) == 0 {
		return action
	}

	if action == "*" {
		action = "[a-zA-Z]+"
	}

	return fmt.Sprintf(`%s\.%s`, prefix, action)
}

