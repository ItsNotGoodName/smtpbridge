package config

import (
	"log"
	"os"
	"path"
	"strconv"

	"github.com/ItsNotGoodName/smtpbridge/domain"
	"github.com/spf13/viper"
)

type Config struct {
	DB        ConfigDB         `json:"database" mapstructure:"database"`
	SMTP      ConfigSMTP       `json:"smtp" mapstructure:"smtp"`
	HTTP      ConfigHTTP       `json:"http" mapstructure:"http"`
	Bridges   []domain.Bridge  `json:"bridges" mapstructure:"bridges"`
	Endpoints []ConfigEndpoint `json:"endpoints" mapstructure:"endpoints"`
}

type ConfigDB struct {
	Type        string `json:"type" mapstructure:"type"`
	DB          string `json:"db" mapstructure:"db"`
	Attachments string `json:"attachments" mapstructure:"attachments"`
}

func (db *ConfigDB) IsBolt() bool {
	return db.Type == "bolt"
}

type ConfigHTTP struct {
	Enable bool   `json:"enable" mapstructure:"enable"`
	Addr   string `json:"-" mapstructure:"-"`
	Host   string `json:"host" mapstructure:"host"`
	Port   uint16 `json:"port" mapstructure:"port"`
}

type ConfigEndpoint struct {
	Name   string            `json:"name" mapstructure:"name"`
	Type   string            `json:"type" mapstructure:"type"`
	Config map[string]string `json:"config" mapstructure:"config"`
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

func New() *Config {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("config.New: could not get user home dir:", err)
	}

	rootPath := path.Join(home, ".smtpbridge")

	return &Config{
		DB: ConfigDB{
			Type:        "bolt",
			DB:          path.Join(rootPath, "smtpbridge.db"),
			Attachments: path.Join(rootPath, "attachments"),
		},
		SMTP: ConfigSMTP{
			Size: 1024 * 1024 * 25,
			Port: 1025,
		},
		HTTP: ConfigHTTP{
			Port: 8080,
		},
		Bridges:   []domain.Bridge{},
		Endpoints: []ConfigEndpoint{},
	}
}

func (c *Config) Load() {
	if err := viper.Unmarshal(c); err != nil {
		log.Fatalln("config.Config.Load:", err)
	}

	c.SMTP.Addr = c.SMTP.Host + ":" + strconv.FormatUint(uint64(c.SMTP.Port), 10)
	c.HTTP.Addr = c.HTTP.Host + ":" + strconv.FormatUint(uint64(c.HTTP.Port), 10)

	if err := os.MkdirAll(c.DB.Attachments, 0755); err != nil {
		log.Fatalln("config.Config.Load: could not create attachments directory", err)
	}

	if err := os.MkdirAll(path.Dir(c.DB.DB), 0755); err != nil {
		log.Fatalln("config.Config.Load: could not create database's parent directory", err)
	}

	log.Printf("config.Config.Load: read %d bridges and %d endpoints", len(c.Bridges), len(c.Endpoints))
}
