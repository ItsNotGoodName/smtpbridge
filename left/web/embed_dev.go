//go:build dev

package web

import (
	"html/template"
	"io/fs"
	"os"
	"path"
)

var projectDIR string

func init() {
	var err error
	projectDIR, err = os.Getwd()
	if err != nil {
		panic(err)
	}
}

func getTemplateFS() fs.FS {
	return os.DirFS(path.Join(projectDIR, packageDIR, templateDIR))
}

func GetAssetFS() fs.FS {
	return os.DirFS(path.Join(projectDIR, packageDIR, assetDIR))
}

func (t *Templater) getTemplate(page string) *template.Template {
	return parseTemplate(getTemplateFS(), page)
}
