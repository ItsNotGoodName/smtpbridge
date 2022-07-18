//go:build !dev

package view

import (
	"embed"
	"html/template"
	"io/fs"
)

//go:embed template
var embedFS embed.FS
var templates map[string]*template.Template

func init() {
	templateFS, err := fs.Sub(embedFS, "template")
	if err != nil {
		panic(err)
	}

	templateFiles, err := fs.ReadDir(templateFS, ".")
	if err != nil {
		panic(err)
	}

	templates = make(map[string]*template.Template)

	for _, file := range templateFiles {
		if file.IsDir() {
			continue
		}

		templates[file.Name()] = parseTemplate(file.Name(), templateFS)
	}
}

func getTemplate(name string) *template.Template {
	return templates[name]
}
