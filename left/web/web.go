package web

import (
	"embed"
	"io/fs"
)

//go:embed dist
var dist embed.FS

func CSSFS() fs.FS {
	f, err := fs.Sub(dist, "dist/css")
	if err != nil {
		panic(err)
	}
	return f
}
