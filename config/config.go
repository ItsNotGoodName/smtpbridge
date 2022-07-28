package config

import (
	"log"
	"os"
	"path"

	"github.com/spf13/viper"
)

type Config struct {
	Memory    bool       `json:"memory" mapstructure:"memory"`
	Directory string     `json:"directory" mapstructure:"directory"`
	Database  Database   `json:"database" mapstructure:"database"`
	Storage   Storage    `json:"storage" mapstructure:"storage"`
	HTTP      HTTP       `json:"http" mapstructure:"http"`
	SMTP      SMTP       `json:"smtp" mapstructure:"smtp"`
	Endpoints []Endpoint `json:"endpoints" mapstructure:"endpoints"`
	Bridges   []Bridge   `json:"bridges" mapstructure:"bridges"`
}

func New() *Config {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("config.New: could not get user's home dir:", err)
	}
	directory := path.Join(home, ".smtpbridge")

	return &Config{
		Directory: directory,
		Database: Database{
			Type: DatabaseTypeBolt,
			Memory: DatabaseMemory{
				Limit: 100,
			},
		},
		Storage: Storage{
			Type: StorageTypeFile,
			Memory: StorageMemory{
				Size: 1024 * 1024 * 100, // 100 MiB
			},
		},
		HTTP: HTTP{
			Port: 8080,
		},
		SMTP: SMTP{
			Size: 1024 * 1024 * 25, // 25 MiB
			Port: 1025,
		},
	}
}

func mustCreatePath(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Printf("config.Config.Load: could not create directory: %s: %s", path, err)
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

	// Set default template for endpoint text
	for i := range c.Endpoints {
		if c.Endpoints[i].TextTemplate == "" {
			c.Endpoints[i].TextTemplate = `FROM: {{ .Message.From }}
SUBJECT: {{ .Message.Subject }}
{{ .Message.Text }}`
		}
	}

	// Override database and storage
	if c.Memory {
		c.Database.Type = DatabaseTypeMemory
		c.Storage.Type = StorageTypeMemory
	}

	// File
	c.Storage.File.Path = path.Join(c.Directory, "attachments")
	if c.Storage.IsFile() {
		mustCreatePath(c.Storage.File.Path)
	}

	// Bolt
	c.Database.Bolt.File = path.Join(c.Directory, "bolt.db")
	if c.Database.IsBolt() {
		mustCreatePath(c.Directory)
	}
}
