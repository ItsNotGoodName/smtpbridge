package web

import (
	"embed"
	"io/fs"
)

var (
	//go:embed dist
	assetsFS embed.FS
)

//go:generate npm run css-build

func GetAssetFS() fs.FS {
	subFS, err := fs.Sub(assetsFS, "dist")
	if err != nil {
		panic(err)
	}
	return subFS
}
