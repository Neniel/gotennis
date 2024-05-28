package config

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type Configuration struct {
	MongoDB          MongoDB          `json:"mongodb"`
	Redis            Redis            `json:"redis"`
	Grafana          Grafana          `json:"grafana"`
	SystemDataSource SystemDataSource `json:"system_data_source"`
}

type SystemDataSource struct {
	URI string `json:"uri"`
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

func ReadFromFile(configFile string) (*Configuration, error) {
	c := Configuration{}
	bs, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(strings.NewReader(string(bs))).Decode(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func ReadFromEnvironmentVariable(configFile string) (*Configuration, error) {
	c := Configuration{}
	err := json.NewDecoder(strings.NewReader(configFile)).Decode(&c)
	if err != nil {
		return nil, err
	}

	return &c, err
}

func LoadConfiguration() (*Configuration, error) {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		return nil, errors.New("environment variable CONFIG_FILE is not set")
	}

	switch os.Getenv("APP_ENVIRONMENT") {
	case "localhost", "docker":
		return ReadFromFile(configFile)
	case "k8s":
		return ReadFromEnvironmentVariable(configFile)
	default:
		return nil, errors.New("invalid value for environment variable APP_ENVIRONMENT")
	}
}
