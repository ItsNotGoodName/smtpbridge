package config

type Database struct {
	Type   string         `json:"type" mapstructure:"type"`
	Memory DatabaseMemory `json:"memory" mapstructure:"memory"`
	Bolt   DatabaseBolt   `json:"bolt" mapstructure:"bolt"`
}

type DatabaseMemory struct {
	Limit int64 `json:"limit" mapstructure:"limit"`
}

type DatabaseBolt struct {
	File string `json:"-" mapstructure:"-"`
}

func (d Database) IsMemory() bool {
	return d.Type == "memory"
}

func (d Database) IsBolt() bool {
	return d.Type == "bolt"
}
