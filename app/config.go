package app

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ConfigSMTP      ConfigSMTP       `json:"smtp" mapstructure:"smtp"`
	Bridges         []Bridge         `json:"bridges" mapstructure:"bridges"`
	ConfigEndpoints []ConfigEndpoint `json:"endpoints" mapstructure:"endpoints"`
}

type ConfigEndpoint struct {
	Name   string            `json:"name" mapstructure:"name"`
	Type   string            `json:"type" mapstructure:"type"`
	Config map[string]string `json:"config" mapstructure:"config"`
}

type ConfigSMTP struct {
	Port string `json:"port" mapstructure:"port"`
	Size int    `json:"size" mapstructure:"size"`
}

func NewConfig() *Config {
	config := &Config{}

	err := viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("app.NewConfig: %s", err)
	}

	log.Printf("app.NewConfig: read %d bridges and %d endpoints", len(config.Bridges), len(config.ConfigEndpoints))

	return config
}
