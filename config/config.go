package config

import (
	"log"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Database ConfigDatabase `json:"database" mapstructure:"database"`
	Storage  ConfigStorage  `json:"storage" mapstructure:"storage"`
	HTTP     ConfigHTTP     `json:"http" mapstructure:"http"`
	SMTP     ConfigSMTP     `json:"smtp" mapstructure:"smtp"`
}

type ConfigDatabase struct {
	Type string `json:"type" mapstructure:"type"`
}

func (cd ConfigDatabase) IsMemDB() bool {
	return cd.Type == "" || cd.Type == "memory"
}

type ConfigStorage struct {
	Type string `json:"type" mapstructure:"type"`
}

func (cs ConfigStorage) IsMemDB() bool {
	return cs.Type == "" || cs.Type == "memory"
}

type ConfigHTTP struct {
	Enable bool   `json:"enable" mapstructure:"enable"`
	Host   string `json:"host" mapstructure:"host"`
	Port   uint16 `json:"port" mapstructure:"port"`
}

func (ch ConfigHTTP) Addr() string {
	return ch.Host + ":" + strconv.FormatUint(uint64(ch.Port), 10)
}

type ConfigSMTP struct {
	Enable   bool   `json:"enable" mapstructure:"enable"`
	Host     string `json:"host" mapstructure:"host"`
	Port     uint16 `json:"port" mapstructure:"port"`
	Size     int    `json:"size" mapstructure:"size"`
	Auth     bool   `json:"auth" mapstructure:"auth"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

func (cs ConfigSMTP) Addr() string {
	return cs.Host + ":" + strconv.FormatUint(uint64(cs.Port), 10)
}

func New() *Config {
	return &Config{
		Database: ConfigDatabase{},
		HTTP: ConfigHTTP{
			Port: 8080,
		},
		SMTP: ConfigSMTP{
			Size: 1024 * 1024 * 25,
			Port: 1025,
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
}
