package utils

import (
	"errors"
	"os"

	"github.com/ametow/getblock/config"

	"gopkg.in/yaml.v3"
)

func ReadConfig(source string) (c *config.Config, err error) {
	raw, err := os.ReadFile(source)
	if err != nil {
		err = errors.New("error reading config from file")
		return
	}
	err = yaml.Unmarshal([]byte(raw), &c)
	if err != nil {
		err = errors.New("error parsing config from yaml")
		c = nil
	}
	return
}
