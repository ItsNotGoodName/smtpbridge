//go:build dev

package asset

import (
	"io/fs"
	"os"
	"path"
	"runtime"
)

var distFS fs.FS

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("no caller information")
	}

	distFS = os.DirFS(path.Join(path.Dir(filename), "dist"))
}

func FS() fs.FS {
	return distFS
}
