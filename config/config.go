package config

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/Neniel/gotennis/log"
)

type Configuration struct {
	MongoDB MongoDB `json:"mongodb"`
	Redis   Redis   `json:"redis"`
	Grafana Grafana `json:"grafana"`
}

type MongoDB struct {
	URI string `json:"uri"`
}

type Redis struct {
	Address  string `json:"address"`
	Password string `json:"password"`
}

type Grafana struct {
	GraphiteToken string `json:"graphite_token"`
}

func LoadConfiguration() (*Configuration, error) {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		log.Logger.Error("environment variable CONFIG_FILE is not set")
		return nil, errors.New("environment variable CONFIG_FILE is not set")
	}

	bs, err := os.ReadFile(configFile)
	if err != nil {
		log.Logger.Error(err.Error())
		return nil, err
	}

	c := Configuration{}
	err = json.NewDecoder(strings.NewReader(string(bs))).Decode(&c)
	if err != nil {
		return nil, err
	}

	return &c, err
}
