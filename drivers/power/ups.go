package power

import (
	"math"

	"github.com/pkg/errors"
	"github.com/timoth-y/chainmetric-sensorsys/drivers/peripheries"
)

// UPSController defines driver for UPS shield device with MAX17040 chip inside.
type UPSController struct {
	*peripheries.I2C
	pwrPin *peripheries.GPIO
}

// NewUPSController constructs new UPSController instance.
func NewUPSController() *UPSController {
	return &UPSController{
		I2C: peripheries.NewI2C(MAX17040_ADDRESS, 1),
		pwrPin: peripheries.NewGPIO(4),
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
	if err := ups.WriteRegBytes(MAX17040_MOD_REG, 0x4000); err != nil {
		return errors.Wrap(err, "failed to set quick-start mode for USP")
	}

	// setup power or reset config
	if err := ups.WriteRegBytes(MAX17040_CMD_REG, 0x0054); err != nil {
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

	return int(math.Round(float64(payload[0]) + float64(payload[1]) / 256)), nil
}

// BatteryVoltage reads current battery voltage.
func (ups *UPSController) BatteryVoltage() (float64, error) {
	payload, err := ups.ReadRegBytes(MAX17040_VOL_REG, 2)
	if err != nil {
		return 0, errors.Wrap(err, "failed to read battery voltage from UPS")
	}

	raw := (payload[0] << 4) | (payload[1] >> 4)

	return float64(raw) * 0.00125, nil
}

// IsPlugged determines whether the UPS is plugged in and charging.
func (ups *UPSController) IsPlugged() bool {
	return ups.pwrPin.IsHigh()
}
