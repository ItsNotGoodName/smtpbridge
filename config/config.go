package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Database  Database   `json:"database" mapstructure:"database"`
	Storage   Storage    `json:"storage" mapstructure:"storage"`
	HTTP      HTTP       `json:"http" mapstructure:"http"`
	SMTP      SMTP       `json:"smtp" mapstructure:"smtp"`
	Endpoints []Endpoint `json:"endpoints" mapstructure:"endpoints"`
}

func New() *Config {
	return &Config{
		Database: Database{
			Memory: DatabaseMemory{
				Limit: 100,
			},
		},
		Storage: Storage{
			Memory: StorageMemory{
				Limit: 30,
				Size:  1024 * 1024 * 100, // 100 MiB
			},
		},
		HTTP: HTTP{
			Enable: true,
			Port:   8080,
		},
		SMTP: SMTP{
			Enable: true,
			Size:   1024 * 1024 * 25, // 25 MiB
			Port:   1025,
		},
	}
}

func (c *Config) Load() {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			log.Fatalln("config.Config.Load:", err)
		}
	}

	if err := viper.Unmarshal(c); err != nil {
		log.Fatalln("config.Config.Load: could not load config:", err)
	}

	if c.HTTP.Disable {
		c.HTTP.Enable = false
	}

	if c.SMTP.Disable {
		c.SMTP.Enable = false
	}
}
