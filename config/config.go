package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Worker  WorkerConfig  `yaml:"worker"`
	Sensors SensorsConfig `yaml:"sensors"`
	Gateway BlockchainConfig `yaml:"gateway"`
}

func ReadConfig(filename string) (sc Config, err error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(file, &sc)
	if err != nil {
		return
	}
	return
}
