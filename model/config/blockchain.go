package config

type BlockchainConfig struct {
	ConnectionConfig string                   `yaml:"connection_config"`
	Identity         BlockchainIdentityConfig `yaml:"identity"`
	ChannelID        string                   `yaml:"channel_id"`
	WalletPath       string                   `yaml:"wallet_path"`
}

type BlockchainIdentityConfig struct {
	Certificate string `yaml:"certificate"`
	PrivateKey  string `yaml:"private_key"`
	UserID      string `yaml:"user_id"`
	OrgID       string `yaml:"org_id"`
	MspID       string `yaml:"msp_id"`
}
