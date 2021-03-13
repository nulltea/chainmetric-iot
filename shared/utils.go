package shared

import "fmt"

func NtoPinName(pin int) string {
	return fmt.Sprintf("GPIO%02d", pin)
}
