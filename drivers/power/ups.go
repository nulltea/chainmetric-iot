package power

import (
	"math"

	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-iot/drivers/periphery"
)

// UPSController defines driver for UPS shield device with MAX17040 chip inside.
type UPSController struct {
	*periphery.I2C
	pwrPin *periphery.GPIO
}

// NewUPSController constructs new UPSController instance.
func NewUPSController() *UPSController {
	return &UPSController{
		I2C:    periphery.NewI2C(MAX17040_ADDRESS, 1),
		pwrPin: periphery.NewGPIO(4),
	}
}

// Init performs initialization sequence of the UPSController.
func (ups *UPSController) Init() error {
	if err := ups.I2C.Init(); err != nil {
		return errors.Wrap(err, "failed to init I2C periphery of UPS")
	}

	if err := ups.pwrPin.Init(); err != nil {
		return errors.Wrap(err, "failed to init GPIO periphery of UPS")
	}

	// enable quick-start mode
	if err := ups.WriteRegBytes(MAX17040_MOD_REG, 0x40, 0x00); err != nil {
		return errors.Wrap(err, "failed to set quick-start mode for USP")
	}

	// reset ups config
	if err := ups.WriteRegBytes(MAX17040_CMD_REG, 0x00, 0x54); err != nil {
		return errors.Wrap(err, "failed to set quick-start mode for USP")
	}

	return nil
}

// BatteryLevel reads current battery level in range of [0-100%].
func (ups *UPSController) BatteryLevel() (int, error) {
	payload, err := ups.ReadRegBytes(MAX17040_SOC_REG, 2)
	if err != nil {
		return 0, errors.Wrap(err, "failed to read battery level from UPS")
	}

	raw := float64(payload[0]) + (float64(payload[1]) / 256)

	if raw > 100 {
		raw = 100
	}

	return int(math.Round(raw)), nil
}

// BatteryVoltage reads current battery voltage.
func (ups *UPSController) BatteryVoltage() (float64, error) {
	payload, err := ups.ReadRegBytes(MAX17040_VOL_REG, 2)
	if err != nil {
		return 0, errors.Wrap(err, "failed to read battery voltage from UPS")
	}

	raw := (payload[0] << 8 + payload[1]) >> 4

	return float64(raw), nil
}

// IsPlugged determines whether the UPS is plugged in and charging.
func (ups *UPSController) IsPlugged() bool {
	return ups.pwrPin.IsHigh()
}
