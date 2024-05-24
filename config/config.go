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

	c := Configuration{}
	appEnvironment := os.Getenv("APP_ENVIRONMENT")
	if appEnvironment == "localhost" || appEnvironment == "docker" {
		bs, err := os.ReadFile(configFile)
		if err != nil {
			log.Logger.Error(err.Error())
			return nil, err
		}

		err = json.NewDecoder(strings.NewReader(string(bs))).Decode(&c)
		if err != nil {
			return nil, err
		}

		return &c, err
	}

	if appEnvironment == "k8s" {
		err := json.NewDecoder(strings.NewReader(configFile)).Decode(&c)
		if err != nil {
			return nil, err
		}

		return &c, err
	}

	return nil, errors.New("invalid value for environment variable APP_ENVIRONMENT")

}
