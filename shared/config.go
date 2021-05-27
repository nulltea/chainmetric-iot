package shared

import (
	"bytes"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// initConfig configures viper from environment variables and configuration files.
func initConfig() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("device.id_file_path", "../device.id")
	viper.SetDefault("device.register_timeout_duration", "1m")
	viper.SetDefault("device.i2c_scan_timeout", "100ms")
	viper.SetDefault("device.hotswap_detect_interval", "3s")
	viper.SetDefault("device.local_cache_path", "/var/cache")
	viper.SetDefault("device.ping_timer_interval", "1m")
	viper.SetDefault("device.battery_check_interval", "1m")

	viper.SetDefault("engine.sensor_sleep_standby_timeout", "1m")

	viper.SetDefault("blockchain.connection_config", "connection.yaml")
	viper.SetDefault("blockchain.identity.certificate", "../identity.pem")
	viper.SetDefault("blockchain.identity.private_key", "../identity.key")
	viper.SetDefault("blockchain.wallet_path", "../keystore")

	viper.SetDefault("bluetooth.enabled", true)
	viper.SetDefault("bluetooth.scan_duration", "1m")
	viper.SetDefault("bluetooth.advertise_duration", "1m")

	viper.SetDefault("sensors.analog.samples_per_read", 100)

	viper.SetDefault("display.enabled", true)
	viper.SetDefault("display.width", 240)
	viper.SetDefault("display.height", 240)
	viper.SetDefault("display.bus", "SPI0.0")
	viper.SetDefault("display.dc_pin", 25)
	viper.SetDefault("display.backlight_pin", 18)
	viper.SetDefault("display.reset_pin", 15)

	viper.SetDefault("mocks.debug_env", false)
	viper.SetDefault("mocks.sensor_duration", "250ms")


	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")

	if err := viper.ReadInConfig(); err != nil {
		Logger.Error(errors.Wrap(err, "failed to read viper config"))
	}
}

// UnmarshalFromConfig retrieves config block by given `key` and decodes it into given structure `v`.
func UnmarshalFromConfig(key string, v interface{}) error {
	bindEnvs(key, v)
	return viper.UnmarshalKey(key, v)
}

// MustUnmarshalFromConfig retrieves config block by given `key` and decodes it into given structure `v`.
// In case of unmarshalling error occurrence it will log fatal error.
func MustUnmarshalFromConfig(key string, v interface{}) {
	if err := UnmarshalFromConfig(key, v); err != nil {
		Logger.Fatal(errors.Wrapf(err, "failed parse config for key '%s'", key))
	}
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
