//go:build dev

package static

import (
	"io/fs"
	"os"
	"path"
	"runtime"
)

var dist fs.FS

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("no caller information")
	}
	packageDir := path.Dir(filename)

	dist = os.DirFS(path.Join(packageDir, "dist"))
}

func CSSFS() fs.FS {
	f, err := fs.Sub(dist, "css")
	if err != nil {
		panic(err)
	}
	return f
}
