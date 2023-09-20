package chiext

import (
	"fmt"
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

func CacheControl(maxAge int) func(h http.Handler) http.Handler {
	var value string
	if maxAge == 0 {
		value = "max-age=31536000,immutable"
	} else {
		value = fmt.Sprintf("max-age=%d", maxAge)
	}
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", value)
			h.ServeHTTP(w, r)
		})
	}
}
