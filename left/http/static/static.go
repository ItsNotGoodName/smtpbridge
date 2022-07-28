//go:build !dev

package static

import (
	"embed"
	"io/fs"
)

//go:embed dist public
var assetFS embed.FS

func AssetFS() fs.FS {
	f, err := fs.Sub(assetFS, "dist")
	if err != nil {
		panic(err)
	}

	return f
}

func FS() fs.FS {
	f, err := fs.Sub(assetFS, "public")
	if err != nil {
		panic(err)
	}

	return f
}
