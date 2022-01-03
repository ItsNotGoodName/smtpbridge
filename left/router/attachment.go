package router

import (
	"io/fs"
	"net/http"
)

func handleImageFS(prefix string, dirFS fs.FS) http.HandlerFunc {
	h := http.StripPrefix(prefix, http.FileServer(http.FS(dirFS)))
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Cache-Control", "max-age=31536000")
		h.ServeHTTP(rw, r)
	}
}
