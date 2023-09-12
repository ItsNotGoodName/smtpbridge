//go:build dev

package web

import (
	"embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

var HeadTags []string

var DevMode bool = true

var FS = mustSubFS(public, "public")

//go:embed public
var public embed.FS

func reloadVite() {
	os.Create(reloadViteFilePath)
}

var reloadViteFilePath string

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to get package directory")
	}

	cwd := filepath.Dir(filename)
	reloadViteFilePath = path.Join(cwd, "reload-vite.local")

	devIP := os.Getenv("DEV_IP")
	if devIP == "" {
		devIP = "127.0.0.1"
	}

	HeadTags = []string{
		fmt.Sprintf(`<script type="module" src="http://%s:5173/@vite/client"></script>`, devIP),
		fmt.Sprintf(`<script type="module" src="http://%s:5173/src/main.ts"></script>`, devIP),
	}
}
