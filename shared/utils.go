package shared

import (
	"fmt"
)

func NtoPinName(pin int) string {
	return fmt.Sprintf("GPIO%02d", pin)
}

func NtoI2cBusName(n int) string {
	return fmt.Sprintf("/dev/i2c-%d", n)
}
