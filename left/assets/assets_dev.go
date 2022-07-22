//go:build dev

package assets

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
	packageDir := path.Dir(filename)

	distFS = os.DirFS(path.Join(packageDir, "dist"))
}

func FS() fs.FS {
	return distFS
}
