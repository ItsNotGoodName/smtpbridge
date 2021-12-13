package app

import (
	"log"

	"github.com/spf13/viper"
)

func NewConfig() *Config {
	// TODO: Make this do less stuff
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("app.NewConfig: %s", err)
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("app.NewConfig: %s", err)
	}

	// TODO: remove this line
	log.Printf("app.NewConfig: loaded %d bridges and %d endpoints\n", len(config.Bridges), len(config.Endpoints))

	return config
}

type Config struct {
	Bridges   []Bridge         `json:"bridges" yaml:"bridges"`
	Endpoints []ConfigEndpoint `json:"endpoints" yaml:"endpoints"`
}

type ConfigEndpoint struct {
	Name   string            `json:"name" yaml:"name"`
	Type   string            `json:"type" yaml:"type"`
	Config map[string]string `json:"config" yaml:"config"`
}
