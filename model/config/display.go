package config

type DisplayConfig struct {
	Width     int   `yaml:"width" mapstructure:"width"`
	Height    int   `yaml:"height" mapstructure:"height"`
	ImageSize int   `yaml:"image_size" mapstructure:"image_size"`
	Rotation  uint8 `yaml:"rotation" mapstructure:"rotation"`
	FrameRate uint8 `yaml:"frame_rate" mapstructure:"frame_rate"`

	Bus          string `yaml:"bus" mapstructure:"bus"`
	DCPin        int    `yaml:"dc_pin" mapstructure:"dc_pin"`
	BacklightPin int    `yaml:"backlight_pin" mapstructure:"backlight_pin"`
	ResetPin     int    `yaml:"reset_pin" mapstructure:"reset_pin"`
}
