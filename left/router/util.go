package router

import (
	"io/fs"
	"net/http"
)

func mwCacheControl(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Cache-Control", "max-age=31536000")
		next(rw, r)
	})
}

func handleFS(prefix string, dirFS fs.FS) http.HandlerFunc {
	h := http.StripPrefix(prefix, http.FileServer(http.FS(dirFS)))
	return func(rw http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(rw, r)
	}
}
