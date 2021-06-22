package utils

import (
	"log"

	"github.com/BurntSushi/toml"
)

//Conf
type Config struct {
	Env testEnv
}

type testEnv struct {
	KeycloakURL      string
	Username         string
	Password         string
	ClientID         string
	BillingServerURL string
	EndpointURL      string
}

func GetConfig(path string) *Config {
	var conf *Config
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		log.Fatal(err)
	}
	return conf
}
