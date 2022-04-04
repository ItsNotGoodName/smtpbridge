package router

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/left/api"
	"github.com/go-chi/chi/v5"
)

const (
	SecondsInDay = 86400
)

// mwCacheControl sets the Cache-Control header to the given value.
func mwCacheControl(maxAge int, next http.HandlerFunc) http.HandlerFunc {
	maxAgeString := fmt.Sprintf("max-age=%d", maxAge)
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Cache-Control", maxAgeString)
		next(rw, r)
	})
}

func handleFS(prefix string, dirFS fs.FS) http.HandlerFunc {
	h := http.StripPrefix(prefix, http.FileServer(http.FS(dirFS)))
	return func(rw http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(rw, r)
	}
}

func render(rd api.Renderer, h api.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rd.Render(w, h(w, r))
	}
}

// mountFS adds GET handlers for all files and folders using the given filesystem.
func mountFS(r chi.Router, f fs.FS) {
	httpFS := http.FS(f)
	fsHandler := http.StripPrefix("/", http.FileServer(httpFS))

	if files, err := fs.ReadDir(f, "."); err == nil {
		for _, f := range files {
			name := f.Name()
			if f.IsDir() {
				r.Get("/"+name+"/*", fsHandler.ServeHTTP)
			} else if name == "index.html" {
				indexHandler := indexGet(httpFS)
				r.Get("/", indexHandler)
				r.Get("/index.html", indexHandler)
			} else {
				r.Get("/"+name, fsHandler.ServeHTTP)
			}
		}
	} else {
		log.Fatal("router.mountFS:", err)
	}
}

// indexGet returns index.html from the given filesystem.
func indexGet(httpFS http.FileSystem) http.HandlerFunc {
	index, err := httpFS.Open("/index.html")
	if err != nil {
		log.Fatal("router.indexGet:", err)
	}

	stat, err := index.Stat()
	if err != nil {
		log.Fatal("router.indexGet:", err)
	}

	modtime := stat.ModTime()

	return func(rw http.ResponseWriter, r *http.Request) {
		http.ServeContent(rw, r, "index.html", modtime, index)
	}
}
