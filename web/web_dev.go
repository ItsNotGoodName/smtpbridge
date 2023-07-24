//go:build dev

package web

import (
	"path"
	"path/filepath"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

const CacheControl = 0

const HeadTags = `<script type="module" src="http://localhost:5173/@vite/client"></script>
<script type="module" src="http://localhost:5173/src/main.ts"></script>`

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

func UseAssets(app *fiber.App) {
	app.Static("/", pathAssets)
}

func Engine() *html.Engine {
	engine := html.New(pathViews, ".html")
	engine.Reload(true)
	return engine
}
