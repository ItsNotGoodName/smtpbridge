//go:build !dev

package web

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
)

var (
	//go:embed template dist
	embedFS    embed.FS
	templateFS fs.FS
	assetFS    fs.FS
)

func init() {
	var err error
	templateFS, err = fs.Sub(embedFS, templateDIR)
	if err != nil {
		panic(err)
	}

	assetFS, err = fs.Sub(embedFS, assetDIR)
	if err != nil {
		panic(err)
	}
}

func (w *Web) getTemplate(name string) *template.Template {
	tmpl, ok := w.templates[(name)]
	if !ok {
		panic(fmt.Errorf("template '%s' not found", name))
	}

	return tmpl
}
