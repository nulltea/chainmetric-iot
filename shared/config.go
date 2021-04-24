package shared

import (
	"bytes"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/timoth-y/iot-blockchain-contracts/shared"
	"gopkg.in/yaml.v2"
)

// InitConfig configures viper from environment variables and configuration files.
func InitConfig() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("device.id_file_path", "../device.id")
	viper.SetDefault("device.hotswap_detect_interval", "3s")

	viper.SetDefault("engine.sensor_sleep_standby_timeout", "1m")

	viper.SetDefault("gateway.connection_config", "connection.yaml")
	viper.SetDefault("gateway.identity.certificate", "../identity.pem")
	viper.SetDefault("gateway.identity.private_key", "../identity.key")
	viper.SetDefault("gateway.wallet_path", "../keystore")

	viper.SetDefault("display.enabled", true)
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
	viper.AddConfigPath("../")

	if err := viper.ReadInConfig(); err != nil {
		shared.Logger.Error(errors.Wrap(err, "failed to read viper config"))
	}
}

func UnmarshalFromConfig(key string, v interface{}) error {
	bindEnvs(key, v)
	return viper.UnmarshalKey(key, v)
}

func bindEnvs(key string, rawVal interface{}) {
	for _, k := range allKeys(key, rawVal) {
		val := viper.Get(k)
		viper.Set(k, val)
	}
}

func allKeys(key string, v interface{}) []string {
	b, err := yaml.Marshal(
		map[string]interface{}{
			key: v,
		},
	)
	if err != nil {
		return nil
	}

	config := viper.New()
	config.SetConfigType("yaml")
	if err := config.ReadConfig(bytes.NewReader(b)); err != nil {
		return nil
	}

	return config.AllKeys()
}
