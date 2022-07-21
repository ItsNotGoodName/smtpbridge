package config

import "strconv"

type SMTP struct {
	Enable   bool   `json:"enable" mapstructure:"enable"`
	Disable  bool   `json:"disable" mapstructure:"disable"`
	Host     string `json:"host" mapstructure:"host"`
	Port     uint16 `json:"port" mapstructure:"port"`
	Size     int    `json:"size" mapstructure:"size"`
	Auth     bool   `json:"auth" mapstructure:"auth"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

func (s SMTP) Addr() string {
	return s.Host + ":" + strconv.FormatUint(uint64(s.Port), 10)
}
