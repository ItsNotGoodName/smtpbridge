package http

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func mwMultiplexAction(get, post, delete http.HandlerFunc) http.HandlerFunc {
	if get == nil {
		get = func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else if post == nil {
		post = func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else if delete == nil {
		delete = func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
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

func mountFS(r chi.Router, f fs.FS) {
	httpFS := http.FS(f)
	fsHandler := http.StripPrefix("/", http.FileServer(httpFS))

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
		log.Fatalln("http.mountFS:", err)
	}
}
