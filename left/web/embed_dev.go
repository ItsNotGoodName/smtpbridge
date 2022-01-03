//go:build dev

package web

import (
	"html/template"
	"io/fs"
	"os"
	"path"
)

func getTemplateFS() fs.FS {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return os.DirFS(path.Join(cwd, packageDIR, templateDIR))
}

func GetAssetFS() fs.FS {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return os.DirFS(path.Join(cwd, packageDIR, assetDIR))
}

func (t *Templater) getTemplate(page string) *template.Template {
	return parseTemplate(getTemplateFS(), page)
}
