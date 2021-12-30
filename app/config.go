package app

import (
	"log"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	AttDir    string           `json:"attachments" mapstructure:"attachments"`
	DBFile    string           `json:"db" mapstructure:"db"`
	SMTP      ConfigSMTP       `json:"smtp" mapstructure:"smtp"`
	HTTP      ConfigHTTP       `json:"http" mapstructure:"http"`
	Bridges   []Bridge         `json:"bridges" mapstructure:"bridges"`
	Endpoints []ConfigEndpoint `json:"endpoints" mapstructure:"endpoints"`
}

type ConfigHTTP struct {
	Enable  bool   `json:"enable" mapstructure:"enable"`
	Address string `json:"-" mapstructure:"-"`
	Host    string `json:"host" mapstructure:"host"`
	Port    uint16 `json:"port" mapstructure:"port"`
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

	config.SMTP.PortStr = strconv.FormatUint(uint64(config.SMTP.Port), 10)
	config.HTTP.Address = config.HTTP.Host + ":" + strconv.FormatUint(uint64(config.HTTP.Port), 10)

	log.Printf("app.NewConfig: read %d bridges and %d endpoints", len(config.Bridges), len(config.Endpoints))

	return config
}
