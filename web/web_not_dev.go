//go:build !dev

package web

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
)

var Development = false

const CacheControl = 3600

//go:embed dist
var assets embed.FS

func AssetsFS() fs.FS {
	f, err := fs.Sub(assets, "dist")
	if err != nil {
		panic(err)
	}
	return f
}

func UseAssets(app *fiber.App) {
	app.Use(filesystem.New(filesystem.Config{
		Root:   http.FS(AssetsFS()),
		MaxAge: CacheControl,
	}))
}

//go:embed views
var views embed.FS

func viewsFS() fs.FS {
	f, err := fs.Sub(views, "views")
	if err != nil {
		panic(err)
	}
	return f
}

func Engine() *html.Engine {
	engine := html.NewFileSystem(http.FS(viewsFS()), ".html")
	return engine
}
