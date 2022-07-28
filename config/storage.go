package config

type Storage struct {
	Type   StorageType   `json:"type" mapstructure:"type"`
	Memory StorageMemory `json:"memory" mapstructure:"memory"`
	File   StorageFile   `json:"file" mapstructure:"file"`
}

type StorageType string

const (
	StorageTypeFile   StorageType = "file"
	StorageTypeMemory StorageType = "memory"
)

func (s Storage) IsFile() bool {
	return s.Type == StorageTypeFile
}

func (s Storage) IsMemory() bool {
	return s.Type == StorageTypeMemory
}

type StorageMemory struct {
	Size int64 `json:"size" mapstructure:"size"`
}

type StorageFile struct {
	Path string `json:"-" mapstructure:"-"`
}
