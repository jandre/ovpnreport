package ovpnreport

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Inputs [](map[string]string) `json: "inputs"`
}

func NewConfig(file string) (*Config, error) {
	var config Config

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
