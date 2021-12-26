package app

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port            string           `json:"port" mapstructure:"port"`
	Size            int              `json:"size" mapstructure:"size"`
	Bridges         []Bridge         `json:"bridges" mapstructure:"bridges"`
	ConfigEndpoints []ConfigEndpoint `json:"endpoints" mapstructure:"endpoints"`
}

type ConfigEndpoint struct {
	Name   string            `json:"name" mapstructure:"name"`
	Type   string            `json:"type" mapstructure:"type"`
	Config map[string]string `json:"config" mapstructure:"config"`
}

func NewConfig() *Config {
	config := &Config{}

	err := viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("app.NewConfig: %s", err)
	}

	// TODO: validate config (e.g. check that all bridges and endpoints have a unique name, make sure all bridges point to valid endpoints, and warn of empty endpoints that are orphaned)

	log.Printf("app.NewConfig: loaded %d bridges and %d endpoints", len(config.Bridges), len(config.ConfigEndpoints))

	return config
}
