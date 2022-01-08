package web

import (
	"html/template"
	"io/fs"
	"log"

	"github.com/ItsNotGoodName/smtpbridge/left"
)

//go:generate npm run css-build

const (
	assetDIR    = "dist"
	templateDIR = "template"
)

// dependencies for pages
var dirs []string = []string{
	"layout/*.html",
	"shim/*.html",
	"message/*.html",
	"attachment/*.html",
}

type Web struct {
	templates map[string]*template.Template
}

func New() *Web {
	return &Web{
		templates: parseTemplates(),
	}
}

func (Web) GetAssetFS() fs.FS {
	return assetFS
}

func (w *Web) GetTemplate(name left.Page) *template.Template {
	return w.getTemplate(string(name))
}

func parseTemplates() map[string]*template.Template {
	tmplFiles, err := fs.ReadDir(templateFS, ".")
	if err != nil {
		log.Fatalln("web.parseTemplate:", err)
	}

	templates := make(map[string]*template.Template)

	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}

		templates[tmpl.Name()] = parseTemplate(tmpl.Name())
	}

	return templates
}

func parseTemplate(name string) *template.Template {
	var patterns []string = append([]string{name}, dirs...)
	pt, err := template.ParseFS(templateFS, patterns...)
	if err != nil {
		panic(err)
	}

	return pt
}
