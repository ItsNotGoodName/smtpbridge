//go:build dev

package view

import (
	"html/template"
	"io/fs"
	"os"
	"path"
	"runtime"
)

var templateFS fs.FS

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("no caller information")
	}

	templateFS = os.DirFS(path.Join(path.Dir(filename), "template"))
}

func getTemplate(name string) *template.Template {
	return parseTemplate(name, templateFS)
}
