package shared

import (
	"os"

	"github.com/op/go-logging"
)

const (
	format = "%{color}%{time:2006.01.02 15:04:05} %{id:04x} %{level}%{color:reset} [%{module}] %{color:bold}%{shortfunc}%{color:reset} -> %{message}"
)

var (
	Logger = logging.MustGetLogger("sensorsys")
)

func InitLogger() {
	backend := logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0),
		logging.MustStringFormatter(format))

	logging.SetBackend(backend)

	level, err := logging.LogLevel(os.Getenv("LOGGING")); if err != nil {
		level = logging.DEBUG
	}

	logging.SetLevel(level, "sensorsys")
}
