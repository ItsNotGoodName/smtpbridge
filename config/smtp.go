package config

import "strconv"

type SMTP struct {
	Disable  bool   `json:"disable" mapstructure:"disable"`
	Host     string `json:"host" mapstructure:"host"`
	Port     uint16 `json:"port" mapstructure:"port"`
	Size     int    `json:"size" mapstructure:"size"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}

func (s SMTP) Addr() string {
	return s.Host + ":" + strconv.FormatUint(uint64(s.Port), 10)
}

func (s SMTP) Auth() bool {
	return s.Username != "" || s.Password != ""
}
