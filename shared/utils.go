package shared

import (
	"encoding/json"
	"fmt"
)

func NtoPinName(pin int) string {
	return fmt.Sprintf("GPIO%02d", pin)
}

func NtoI2cBusName(n int) string {
	return fmt.Sprintf("/dev/i2c-%d", n)
}

func PrettyPrint(obj interface{}) string {
	pretty, err := json.MarshalIndent(obj, "", "\t"); if err != nil {
		return err.Error()
	}
	return string(pretty)
}
