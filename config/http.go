package config

import "strconv"

type HTTP struct {
	Enable  bool   `json:"enable" mapstructure:"enable"`
	Disable bool   `json:"disable" mapstructure:"disable"`
	Host    string `json:"host" mapstructure:"host"`
	Port    uint16 `json:"port" mapstructure:"port"`
}

func (h HTTP) Addr() string {
	return h.Host + ":" + strconv.FormatUint(uint64(h.Port), 10)
}
