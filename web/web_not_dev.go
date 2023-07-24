//go:build !dev

package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
)

const CacheControl = 3600

var HeadTags = ""

//go:embed dist
var assets embed.FS

func assetsFS() fs.FS {
	f, err := fs.Sub(assets, "dist")
	if err != nil {
		panic(err)
	}
	return f
}

func UseAssets(app *fiber.App) {
	app.Use(filesystem.New(filesystem.Config{
		Root:   http.FS(assetsFS()),
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

func init() {
	type Manifest struct {
		CSS     []string `json:"css"`
		File    string   `json:"file"`
		IsEntry bool     `json:"isEntry"`
		Src     string   `json:"src"`
	}

	file, err := assets.Open("dist/manifest.json")
	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var manifestMap map[string]Manifest
	if err := json.Unmarshal(data, &manifestMap); err != nil {
		panic(err)
	}

	var manifest Manifest
	func() {
		for _, man := range manifestMap {
			if man.IsEntry {
				manifest = man
				return
			}
		}
		panic("entrypoint not found in manifest")
	}()

	for _, v := range manifest.CSS {
		HeadTags += fmt.Sprintf(`<link rel="stylesheet" href="/%s" />`, v)
	}
	HeadTags += fmt.Sprintf(`<script type="module" src="/%s"></script>`, manifest.File)
}
