package helpers

import (
	"encoding/json"
	"html/template"
	"io"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/build"
	"github.com/ItsNotGoodName/smtpbridge/web"
	"github.com/dustin/go-humanize"
)

var Map template.FuncMap = template.FuncMap{
	"build": func() build.Build {
		return build.Current
	},
	"development": func() bool {
		return web.Development
	},
	"manifest": func() _Manifest {
		return manifest
	},
	"timeFormat": func(date time.Time) string {
		return date.Format("Jan _2 2006 15:04:05")
	},
	"timeHumanize": func(date time.Time) string {
		return humanize.Time(date)
	},
	"bytesHumanize": func(bytes int64) string {
		return humanize.Bytes(uint64(bytes))
	},
	"json": func(data any) string {
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			panic(err)
		}

		return string(jsonData)
	},
	"unescape": func(s string) template.HTML {
		return template.HTML(s)
	},
	"query": func(queries map[string]string, vals ...any) template.URL {
		return template.URL(Query(queries, vals...))
	},
}

var manifest _Manifest

type _Manifest struct {
	CSS     []string `json:"css"`
	File    string   `json:"file"`
	IsEntry bool     `json:"isEntry"`
	Src     string   `json:"src"`
}

func init() {
	fs := web.AssetsFS()
	if web.Development {
		return
	}

	file, err := fs.Open("manifest.json")
	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var manifestMap map[string]_Manifest
	if err := json.Unmarshal(data, &manifestMap); err != nil {
		panic(err)
	}
	mustSetManifest(manifestMap)
}

func mustSetManifest(manifestMap map[string]_Manifest) {
	for _, man := range manifestMap {
		if man.IsEntry {
			manifest = man
			return
		}
	}
	panic("entrypoint not found in manifest")
}
