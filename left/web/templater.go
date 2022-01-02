package web

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type Page string

const (
	PageIndex = "index.html"
)

const (
	shimDIR     = "template/shim"
	templateDIR = "template"
	templateEXT = "/*.html"
)

var (
	//go:embed template
	templateFS embed.FS
)

type Templater struct {
	templates map[string]*template.Template
}

func NewTemplater() *Templater {
	t := Templater{
		templates: make(map[string]*template.Template),
	}

	tmplFiles, err := fs.ReadDir(templateFS, templateDIR)
	if err != nil {
		log.Fatalln("web.NewTemplater:", err)
	}

	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}

		pt, err := template.ParseFS(templateFS, templateDIR+"/"+tmpl.Name(), shimDIR+templateEXT)
		if err != nil {
			log.Fatalln("web.NewTemplater:", err)
		}

		t.templates[tmpl.Name()] = pt
	}

	return &t
}

func (t *Templater) Render(page Page, rw http.ResponseWriter, data interface{}) {
	err := t.templates[string(page)].Execute(rw, data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
