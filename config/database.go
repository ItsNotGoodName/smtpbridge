package config

type Database struct {
	Type   DatabaseType   `json:"type" mapstructure:"type"`
	Memory DatabaseMemory `json:"memory" mapstructure:"memory"`
	Bolt   DatabaseBolt   `json:"bolt" mapstructure:"bolt"`
}

type DatabaseType string

const (
	DatabaseTypeBolt   = "bolt"
	DatabaseTypeMemory = "memory"
)

type DatabaseMemory struct {
	Limit int64 `json:"limit" mapstructure:"limit"`
}

type DatabaseBolt struct {
	File string `json:"-" mapstructure:"-"`
}

func (d Database) IsMemory() bool {
	return d.Type == DatabaseTypeMemory
}

func (d Database) IsBolt() bool {
	return d.Type == DatabaseTypeBolt
}
