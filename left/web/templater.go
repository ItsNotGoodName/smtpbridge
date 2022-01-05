package web

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path"
)

//go:generate npm run css-build

type Page string

const (
	PageIndex   Page = "index.html"
	PageMessage Page = "message.html"
	PageInfo    Page = "info.html"
)

const (
	packageDIR  = "left/web"
	assetDIR    = "dist"
	templateDIR = "template"
	templateEXT = "/*.html"
	shimDIR     = "shim"
)

type Templater struct {
	templates map[string]*template.Template
}

func NewTemplater() *Templater {
	return &Templater{
		templates: parseTemplateFS(getTemplateFS()),
	}
}

func (t *Templater) Render(page Page, rw http.ResponseWriter, data interface{}) {
	err := t.getTemplate(string(page)).Execute(rw, data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func parseTemplateFS(dirFS fs.FS) map[string]*template.Template {
	tmplFiles, err := fs.ReadDir(dirFS, ".")
	if err != nil {
		log.Fatalln("web.parseTemplateFS:", err)
	}

	templates := make(map[string]*template.Template)

	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}

		templates[tmpl.Name()] = parseTemplate(dirFS, tmpl.Name())
	}

	return templates
}

func parseTemplate(dirFS fs.FS, name string) *template.Template {
	pt, err := template.ParseFS(dirFS, name, path.Join(shimDIR, templateEXT))
	if err != nil {
		panic(err)
	}

	return pt
}
