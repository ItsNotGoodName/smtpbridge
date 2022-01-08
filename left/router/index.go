package router

import (
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/ItsNotGoodName/smtpbridge/left"
	"github.com/ItsNotGoodName/smtpbridge/pkg/paginate"
)

func handleIndexGet(w left.WebRepository, a *app.App) http.HandlerFunc {
	param := "page"

	return func(rw http.ResponseWriter, r *http.Request) {
		var page int
		if p := r.URL.Query().Get(param); p != "" {
			if i, err := strconv.Atoi(p); err == nil {
				page = i
			}
		}

		res, err := a.MessageList(&app.MessageListRequest{Page: page})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		data := left.IndexData{
			Messages: res.Messages,
			Paginate: paginate.New(*r.URL, param, res.PageMin, res.Page, res.PageMax),
		}
		render(rw, w, left.IndexPage, &data)
	}
}
