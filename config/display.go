package config

type DisplayConfig struct {
	Width     int   `yaml:"width"`
	Height    int   `yaml:"height"`
	ImageSize int   `yaml:"imageSize"`
	Rotation  uint8 `yaml:"rotation"`
	FrameRate uint8 `yaml:"frameRate"`

	Bus          string `yaml:"bus"`
	DCPin        int    `yaml:"dcPin"`
	BacklightPin int    `yaml:"backlightPin"`
	ResetPin     int    `yaml:"resetPin"`
}
