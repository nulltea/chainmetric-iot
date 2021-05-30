package shared

import (
	"os"

	"github.com/op/go-logging"
)

const (
	format = "%{color}%{time:2006.01.02 15:04:05} %{id:03x} %{level:-5s}%{color:reset} [%{module}] %{color:bold}%{shortfunc}%{color:reset} -> %{message}"
)

// Logger is an instance of the shared logger tool.
var Logger = logging.MustGetLogger("sensorsys")

func initLogger() {
	backend := logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0),
		logging.MustStringFormatter(format))

	logging.SetBackend(backend)

	level, err := logging.LogLevel(os.Getenv("LOGGING")); if err != nil {
		level = logging.DEBUG
	}

	logging.SetLevel(level, "sensorsys")
}
