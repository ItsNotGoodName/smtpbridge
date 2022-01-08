package router

import (
	"io/fs"
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/left"
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

func render(rw http.ResponseWriter, w left.WebRepository, page left.Page, data interface{}) {
	err := w.GetTemplate(page).Execute(rw, data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
