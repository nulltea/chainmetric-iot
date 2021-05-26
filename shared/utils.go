package shared

import (
	"fmt"

	"github.com/pkg/errors"
)

// NtoPinName returns GPIO pin name based on specified `pin` number.
func NtoPinName(pin int) string {
	return fmt.Sprintf("GPIO%d", pin)
}

// NtoI2cBusName returns I2C bus name based on specified `n` number.
func NtoI2cBusName(n int) string {
	return fmt.Sprintf("/dev/i2c-%d", n)
}

// MustExecute executes `fn` function and in case of error logs it, followed by a call to os.Exit(1).
// Use `msg` to specify details of error to wrap by.
func MustExecute(fn func() error, msg string) {
	if err := fn(); err != nil {
		Logger.Fatal(errors.Wrap(err, msg))
	}
}

// Execute executes `fn` function and in case of error logs it.
// Use `msg` to specify details of error to wrap by.
func Execute(fn func() error, msg string) {
	if err := fn(); err != nil {
		Logger.Error(errors.Wrap(err, msg))
	}
}
