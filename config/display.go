package config

type DisplayConfig struct {
	Width int `yaml:"width"`
	Height int `yaml:"height"`
	Rotation uint8 `yaml:"rotation"`
	FrameRate uint8 `yaml:"frameRate"`
}
