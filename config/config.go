package config

import (
	"log"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Database ConfigDatabase `json:"database" mapstructure:"database"`
	SMTP     ConfigSMTP     `json:"smtp" mapstructure:"smtp"`
	HTTP     ConfigHTTP     `json:"http" mapstructure:"http"`
}

type ConfigDatabase struct {
	Type string `json:"type" mapstructure:"type"`
}

func (db ConfigDatabase) IsMemDB() bool {
	return db.Type == ""
}

type ConfigSMTP struct {
	Addr     string `json:"-" mapstructure:"-"`
	Host     string `json:"host" mapstructure:"host"`
	Port     uint16 `json:"port" mapstructure:"port"`
	Size     int    `json:"size" mapstructure:"size"`
	Auth     bool   `json:"auth" mapstructure:"auth"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

type ConfigHTTP struct {
	Enable bool   `json:"enable" mapstructure:"enable"`
	Addr   string `json:"-" mapstructure:"-"`
	Host   string `json:"host" mapstructure:"host"`
	Port   uint16 `json:"port" mapstructure:"port"`
}

func New() *Config {
	return &Config{
		Database: ConfigDatabase{},
		SMTP: ConfigSMTP{
			Size: 1024 * 1024 * 25,
			Port: 1025,
		},
		HTTP: ConfigHTTP{
			Port: 8080,
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

	c.SMTP.Addr = c.SMTP.Host + ":" + strconv.FormatUint(uint64(c.SMTP.Port), 10)
	c.HTTP.Addr = c.HTTP.Host + ":" + strconv.FormatUint(uint64(c.HTTP.Port), 10)
}
