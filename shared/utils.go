package shared

import (
	"fmt"

	"github.com/pkg/errors"
)

func NtoPinName(pin int) string {
	return fmt.Sprintf("GPIO%d", pin)
}

func NtoI2cBusName(n int) string {
	return fmt.Sprintf("/dev/i2c-%d", n)
}

func MustExecute(fn func() error, msg string) {
	if err := fn(); err != nil {
		Logger.Fatal(errors.Wrap(err, msg))
	}
}
