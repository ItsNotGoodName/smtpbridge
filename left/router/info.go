package router

import (
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/ItsNotGoodName/smtpbridge/left"
)

func handleInfoGet(w left.WebRepository, a *app.App) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		info, err := a.Info()
		if err != nil {
			renderError(rw, err, http.StatusInternalServerError)
			return
		}

		data := left.InfoData{Info: *info}
		render(rw, w, left.InfoPage, &data)
	}
}
