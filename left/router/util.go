package router

import (
	"fmt"
	"io/fs"
	"net/http"
)

func mwMultiplexAction(get, post, delete http.HandlerFunc) http.HandlerFunc {
	if get == nil {
		get = func(w http.ResponseWriter, r *http.Request) {}
	} else if post == nil {
		post = func(w http.ResponseWriter, r *http.Request) {}
	} else if delete == nil {
		delete = func(w http.ResponseWriter, r *http.Request) {}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		action := r.URL.Query().Get("action")
		if action == "delete" {
			delete(w, r)
		} else if action == "post" {
			post(w, r)
		} else {
			get(w, r)
		}
	}
}

func mwCacheControl(next http.HandlerFunc, maxAge int) http.HandlerFunc {
	maxAgeString := fmt.Sprintf("max-age=%d", maxAge)
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Cache-Control", maxAgeString)
		next(rw, r)
	})
}

func handlePrefixFS(prefix string, fs fs.FS) http.HandlerFunc {
	h := http.StripPrefix(prefix, http.FileServer(http.FS(fs)))
	return func(rw http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(rw, r)
	}
}
