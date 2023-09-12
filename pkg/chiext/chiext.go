package chiext

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// MountFS adds GET handlers for all files and directories for the given fs.FS.
func MountFS(r chi.Router, f fs.FS) {
	fsHandler := http.StripPrefix("/", http.FileServer(http.FS(f)))

	if files, err := fs.ReadDir(f, "."); err == nil {
		for _, f := range files {
			name := f.Name()
			if f.IsDir() {
				r.Get("/"+name+"/*", fsHandler.ServeHTTP)
			} else {
				r.Get("/"+name, fsHandler.ServeHTTP)
			}
		}
	} else if err != fs.ErrNotExist {
		panic(err)
	}
}
