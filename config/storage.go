package config

type Storage struct {
	Type   string        `json:"type" mapstructure:"type"`
	Memory StorageMemory `json:"memory" mapstructure:"memory"`
}

type StorageMemory struct {
	Limit int64 `json:"limit" mapstructure:"limit"`
	Size  int64 `json:"size" mapstructure:"size"`
}

func (s Storage) IsMemDB() bool {
	return s.Type == "" || s.Type == "memory"
}
