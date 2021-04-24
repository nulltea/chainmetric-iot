package shared

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/timoth-y/iot-blockchain-contracts/shared"
)

// InitConfig configures viper from environment variables and configuration files.
func InitConfig() {
	viper.AutomaticEnv()

	viper.SetDefault("engine.sensor_sleep_standby_timeout", "1m")

	viper.SetDefault("gateway.connection_config", "connection.yaml")
	viper.SetDefault("gateway.identity.certificate", "../identity.pem")
	viper.SetDefault("gateway.identity.private_key", "../identity.key")
	viper.SetDefault("gateway.wallet_path", "../keystore")

	viper.SetDefault("display.width", 240)
	viper.SetDefault("display.height", 240)
	viper.SetDefault("display.image_size", 150)
	viper.SetDefault("display.bus", "SPI0.0")
	viper.SetDefault("display.dc_pin", 25)
	viper.SetDefault("display.backlight_pin", 18)
	viper.SetDefault("display.reset_pin", 15)

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		shared.Logger.Error(errors.Wrap(err, "failed to read viper config"))
	}
}
