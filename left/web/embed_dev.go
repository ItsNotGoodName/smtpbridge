//go:build dev

package web

import (
	"html/template"
	"io/fs"
	"os"
	"path"
	"runtime"
)

var (
	templateFS fs.FS
	assetFS    fs.FS
)

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("no caller information")
	}
	packageDIR := path.Dir(filename)

	templateFS = os.DirFS(path.Join(packageDIR, templateDIR))
	assetFS = os.DirFS(path.Join(packageDIR, assetDIR))
}

func (Web) getTemplate(page string) *template.Template {
	return parseTemplate(page)
}
