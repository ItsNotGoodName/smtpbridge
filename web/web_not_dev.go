//go:build !dev

package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
)

var HeadTags []string

var DevMode bool = false

var FS fs.FS

//go:embed dist
var dist embed.FS

func reloadVite() {}

func init() {
	FS = mustSubFS(dist, "dist")

	// Parse manifest
	file, err := FS.Open("manifest.json")
	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	type Manifest struct {
		CSS     []string `json:"css"`
		File    string   `json:"file"`
		IsEntry bool     `json:"isEntry"`
		Src     string   `json:"src"`
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
		HeadTags = append(HeadTags, fmt.Sprintf(`<link rel="stylesheet" href="/%s" />`, v))
	}
	HeadTags = append(HeadTags, fmt.Sprintf(`<script type="module" src="/%s"></script>`, manifest.File))
}
