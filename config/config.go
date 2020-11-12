package config

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type Config struct {
	Username              string
	Password              string
	ReturnRoute           string
	ReturnRtransportation string
	LeaveRoute            string
	LeaveTransportation   string
	Location              string
	ParentsPhone          string
	LeaveReason           string
	Contact               string
	SCKey                 string
}

func (c *Config) Get() error {
	tomlFile, err := ioutil.ReadFile("config.toml")
	if err != nil {
		return err
	}
	if _, err = toml.Decode(string(tomlFile), c); err != nil {
		return err
	}

	return nil
}
