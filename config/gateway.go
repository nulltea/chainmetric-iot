package config

type BlockchainConfig struct {
	ConnectionConfig string                   `yaml:"connectionConfig"`
	Identity         BlockchainIdentityConfig `yaml:"identity"`
	ChannelID        string                   `yaml:"channelID"`
	WalletPath       string                   `yaml:"walletPath"`
}

type BlockchainIdentityConfig struct {
	Certificate string `yaml:"certificate"`
	PrivateKey  string `yaml:"privateKey"`
	UserID      string `yaml:"userID"`
	OrgID       string `yaml:"orgID"`
	MspID       string `yaml:"mspID"`
}
