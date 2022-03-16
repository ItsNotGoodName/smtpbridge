package web

import (
	"embed"
	"io/fs"
	"log"
	"mime"
)

func init() {
	mime.AddExtensionType(".js", "application/javascript")
}

//go:generate npm run build

//go:embed dist
var dist embed.FS

func FS() fs.FS {
	f, err := fs.Sub(dist, "dist")
	if err != nil {
		log.Fatal("web.FS:", err)
	}
	return f
}
