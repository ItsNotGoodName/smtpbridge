package config

type Endpoint struct {
	Name   string            `json:"name" mapstructure:"name"`
	Type   string            `json:"type" mapstructure:"type"`
	Config map[string]string `json:"config" mapstructure:"config"`
}
