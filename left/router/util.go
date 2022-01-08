package router

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/left"
)

const (
	SecondsInYear = 31536000
	SecondsInDay  = 86400
)

func mwCacheControl(next http.HandlerFunc, maxAge int) http.HandlerFunc {
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

func render(rw http.ResponseWriter, w left.WebRepository, page left.Page, data interface{}) {
	err := w.GetTemplate(page).Execute(rw, data)
	if err != nil {
		renderError(rw, err, http.StatusInternalServerError)
	}
}

func renderError(rw http.ResponseWriter, err error, status int) {
	http.Error(rw, err.Error(), status)
}
