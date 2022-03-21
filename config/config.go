package config

import (
	"log"
	"os"
	"path"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Storage   ConfigStorage    `json:"storage" mapstructure:"storage"`
	Database  ConfigDatabase   `json:"database" mapstructure:"database"`
	HTTP      ConfigHTTP       `json:"http" mapstructure:"http"`
	SMTP      ConfigSMTP       `json:"smtp" mapstructure:"smtp"`
	Bridges   []ConfigBridge   `json:"bridges" mapstructure:"bridges"`
	Endpoints []ConfigEndpoint `json:"endpoints" mapstructure:"endpoints"`
}

type ConfigBridge struct {
	Name          string                 `json:"name" mapstructure:"name"`
	NoText        bool                   `json:"no_text" mapstructure:"no_text"`
	NoAttachments bool                   `json:"no_attachments" mapstructure:"no_attachments"`
	Endpoints     []ConfigBridgeEndpoint `json:"endpoints" mapstructure:"endpoints"`
	Filters       []ConfigFilter         `json:"filters" mapstructure:"filters"`
}

type ConfigBridgeEndpoint struct {
	Name             string `json:"name" mapstructure:"name"`
	NoTextStr        string `json:"no_text" mapstructure:"no_text"`
	NoAttachmentsStr string `json:"no_attachments" mapstructure:"no_attachments"`
	NoText           bool   `json:"-" mapstructure:"-"`
	NoAttachments    bool   `json:"-" mapstructure:"-"`
}

type ConfigFilter struct {
	To        string `json:"to,omitempty" mapstructure:"to,omitempty"`
	From      string `json:"from,omitempty" mapstructure:"from,omitempty"`
	ToRegex   string `json:"to_regex,omitempty" mapstructure:"to_regex,omitempty"`
	FromRegex string `json:"from_regex,omitempty" mapstructure:"from_regex,omitempty"`
}

type ConfigDatabase struct {
	Type     string `json:"type" mapstructure:"type"`
	BoltFile string `json:"-" mapstructure:"-"`
	BoltPath string `json:"-" mapstructure:"-"`
}

type ConfigStorage struct {
	Size                 int64  `json:"size" mapstructure:"size"`
	Path                 string `json:"path" mapstructure:"path"`
	AttachmentsDirectory string `json:"-" mapstructure:"-"`
	AttachmentsPath      string `json:"-" mapstructure:"-"`
}

func (db ConfigDatabase) IsBolt() bool {
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
		log.Fatalln("config.New: could not get user's home dir:", err)
	}

	return &Config{
		Storage: ConfigStorage{
			Size:                 1024 * 1024 * 2048,
			Path:                 path.Join(home, ".smtpbridge"),
			AttachmentsDirectory: "attachments",
		},
		Database: ConfigDatabase{
			BoltFile: "bolt.db",
		},
		HTTP: ConfigHTTP{
			Port: 8080,
		},
		SMTP: ConfigSMTP{
			Size: 1024 * 1024 * 25,
			Port: 1025,
		},
		Bridges:   []ConfigBridge{},
		Endpoints: []ConfigEndpoint{},
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
	c.Storage.AttachmentsPath = path.Join(c.Storage.Path, c.Storage.AttachmentsDirectory)
	c.Database.BoltPath = path.Join(c.Storage.Path, c.Database.BoltFile)

	for _, bridge := range c.Bridges {
		for j, endpoint := range bridge.Endpoints {
			if endpoint.NoTextStr != "" {
				bridge.Endpoints[j].NoText, _ = strconv.ParseBool(endpoint.NoTextStr)
			} else {
				bridge.Endpoints[j].NoText = bridge.NoText
			}
			if endpoint.NoAttachmentsStr != "" {
				bridge.Endpoints[j].NoAttachments, _ = strconv.ParseBool(endpoint.NoAttachmentsStr)
			} else {
				bridge.Endpoints[j].NoAttachments = bridge.NoAttachments
			}
		}
	}

	if err := os.MkdirAll(c.Storage.Path, 0755); err != nil {
		log.Println("config.Config.Load: could not create storage directory:", err)
	}

	if err := os.MkdirAll(c.Storage.AttachmentsPath, 0755); err != nil {
		log.Println("config.Config.Load: could not create attachments directory:", err)
	}

	log.Printf("config.Config.Load: %d bridges and %d endpoints", len(c.Bridges), len(c.Endpoints))
}
