package shared

import (
	"github.com/spf13/viper"
	"github.com/timoth-y/go-eventdriver"
)

// InitCore performs core dependencies initialization sequence.
func InitCore() {
	initLogger()
	initConfig()
	initLevelDB()
	initPeriphery()

	eventdriver.Init(
		eventdriver.WithLogger(Logger),
		eventdriver.WithBufferSize(viper.GetInt("local_events_bugger_size")),
	)
}

// CloseCore performs core dependencies close sequence.
func CloseCore() {
	closeLevelDB()
	closePeriphery()

	eventdriver.Close()
}
