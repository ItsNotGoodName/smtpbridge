package config

type Database struct {
	Type   string         `json:"type" mapstructure:"type"`
	Memory DatabaseMemory `json:"memory" mapstructure:"memory"`
}

type DatabaseMemory struct {
	Limit int64 `json:"limit" mapstructure:"limit"`
}

func (d Database) IsMemDB() bool {
	return d.Type == "" || d.Type == "memory"
}
