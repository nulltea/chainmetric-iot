package config

type BlockchainConfig struct {
	ConnectionConfig string                   `yaml:"connection_config" mapstructure:"connection_config"`
	Identity         BlockchainIdentityConfig `yaml:"identity" mapstructure:"identity"`
	ChannelID        string                   `yaml:"channel_id" mapstructure:"channel_id"`
	WalletPath       string                   `yaml:"wallet_path" mapstructure:"wallet_path"`
}

type BlockchainIdentityConfig struct {
	Certificate string `yaml:"certificate" mapstructure:"certificate"`
	PrivateKey  string `yaml:"private_key" mapstructure:"private_key"`
	UserID      string `yaml:"user_id" mapstructure:"user_id"`
	OrgID       string `yaml:"org_id" mapstructure:"org_id"`
	MspID       string `yaml:"msp_id" mapstructure:"msp_id"`
}
