package web

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed dist
var dist embed.FS

func AssetFS() fs.FS {
	f, err := fs.Sub(dist, "dist")
	if err != nil {
		log.Fatalln("web.AssetFS:", err)
	}
	return f
}
