package conf

import (
	"io/ioutil"
	"os"

	"github.com/somedevv/permit-ssh/colors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DB string `yaml:"db type"`
}

func (c *Config) GetConf() *Config {

	confFile, err := ioutil.ReadFile(os.Getenv("HOME") + "/.local/bin/.permit_data/config.yaml")
	if err != nil {
		colors.Red.Println("Could not load config file, using default values")
		return c.setDefaultConf()
	}
	err = yaml.Unmarshal(confFile, c)
	if err != nil {
		colors.Red.Println("Unmarshal: %v", err)
		os.Exit(1)
	}

	return c
}

func (c *Config) setDefaultConf() *Config {

	if c.DB == "" {
		c.DB = "local"
	}

	return c
}
