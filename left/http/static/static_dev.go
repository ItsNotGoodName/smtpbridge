//go:build dev

package static

import (
	"io/fs"
	"os"
	"path"
	"runtime"
)

var distFS fs.FS
var publicFS fs.FS

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("no caller information")
	}

	distFS = os.DirFS(path.Join(path.Dir(filename), "dist"))
	publicFS = os.DirFS(path.Join(path.Dir(filename), "public"))
}

func AssetFS() fs.FS {
	return distFS
}

func FS() fs.FS {
	return publicFS
}
