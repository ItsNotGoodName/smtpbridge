package app

import (
	"log"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	ConfigSMTP      ConfigSMTP       `json:"smtp" mapstructure:"smtp"`
	DB              string           `json:"db" mapstructure:"db"`
	Bridges         []Bridge         `json:"bridges" mapstructure:"bridges"`
	ConfigEndpoints []ConfigEndpoint `json:"endpoints" mapstructure:"endpoints"`
}

type ConfigEndpoint struct {
	Name   string            `json:"name" mapstructure:"name"`
	Type   string            `json:"type" mapstructure:"type"`
	Config map[string]string `json:"config" mapstructure:"config"`
}

type ConfigSMTP struct {
	Host    string `json:"host" mapstructure:"host"`
	Port    uint16 `json:"port" mapstructure:"port"`
	PortStr string `json:"-" mapstructure:"-"`
	Size    int    `json:"size" mapstructure:"size"`
}

func NewConfig() *Config {
	config := &Config{}

	err := viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("app.NewConfig: %s", err)
	}

	config.ConfigSMTP.PortStr = strconv.FormatUint(uint64(config.ConfigSMTP.Port), 10)

	log.Printf("app.NewConfig: read %d bridges and %d endpoints", len(config.Bridges), len(config.ConfigEndpoints))

	return config
}
