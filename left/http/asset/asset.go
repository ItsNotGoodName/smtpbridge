//go:build !dev

package asset

import (
	"embed"
	"io/fs"
)

//go:embed dist
var assetFS embed.FS

func FS() fs.FS {
	f, err := fs.Sub(assetFS, "dist")
	if err != nil {
		panic(err)
	}

	return f
}
