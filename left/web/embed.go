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
	embedFS embed.FS
)

func getTemplateFS() fs.FS {
	subFS, err := fs.Sub(embedFS, templateDIR)
	if err != nil {
		panic(err)
	}
	return subFS
}

func GetAssetFS() fs.FS {
	subFS, err := fs.Sub(embedFS, assetDIR)
	if err != nil {
		panic(err)
	}

	return subFS
}

func (t *Templater) getTemplate(name string) *template.Template {
	tmpl, ok := t.templates[name]
	if !ok {
		panic(fmt.Errorf("template '%s' not found", name))
	}

	return tmpl
}
