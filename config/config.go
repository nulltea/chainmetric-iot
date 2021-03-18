package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/timoth-y/iot-blockchain-sensorsys/shared"
)

type Config struct {
	Worker  WorkerConfig     `yaml:"worker"`
	Sensors SensorsConfig    `yaml:"sensors"`
	Display DisplayConfig    `yaml:"display"`
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

func MustReadConfig(filename string) Config {
	cnf, err := ReadConfig(filename); if err != nil {
		shared.Logger.Fatal(cnf)
	}

	return cnf
}
