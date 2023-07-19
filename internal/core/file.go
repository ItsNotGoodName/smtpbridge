package core

import (
	"io/fs"
	"os"
)

type FileStore struct {
	Dir string
	FS  fs.FS
}

func NewFileStore(dir string) FileStore {
	return FileStore{
		Dir: dir,
		FS:  os.DirFS(dir),
	}
}
