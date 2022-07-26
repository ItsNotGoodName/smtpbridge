package config

type Storage struct {
	Type      string           `json:"type" mapstructure:"type"`
	Memory    StorageMemory    `json:"memory" mapstructure:"memory"`
	Directory StorageDirectory `json:"directory" mapstructure:"directory"`
}

func (s Storage) IsMemory() bool {
	return s.Type == "memory"
}

func (s Storage) IsDirectory() bool {
	return s.Type == "directory"
}

type StorageMemory struct {
	Size int64 `json:"size" mapstructure:"size"`
}

type StorageDirectory struct {
	Path string `json:"-" mapstructure:"-"`
}
