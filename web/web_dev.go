//go:build dev

package web

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

var Development = true

const CacheControl = 0

var pathAssets string
var pathViews string

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to get package directory")
	}

	cwd := filepath.Dir(filename)
	pathAssets = path.Join(cwd, "public")
	pathViews = path.Join(cwd, "views")
}

func AssetsFS() fs.FS {
	return os.DirFS(pathAssets)
}

func UseAssets(app *fiber.App) {
	app.Static("/", pathAssets)
}

func Engine() *html.Engine {
	engine := html.New(pathViews, ".html")
	engine.Reload(true)
	return engine
}
