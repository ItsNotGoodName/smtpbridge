package router

import (
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/ItsNotGoodName/smtpbridge/left/web"
)

func handleInfoGet(t *web.Templater, a *app.App) http.HandlerFunc {
	type Data struct {
		Info *app.InfoResponse
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		info, err := a.Info()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		t.Render(web.PageInfo, rw, &Data{
			Info: info,
		})
	}
}
