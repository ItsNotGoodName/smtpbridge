package router

import (
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/ItsNotGoodName/smtpbridge/left"
	"github.com/ItsNotGoodName/smtpbridge/pkg/paginate"
)

func handleIndexGet(w left.WebRepository, a *app.App) http.HandlerFunc {
	pageParam := "page"
	statusParam := "status"

	return func(rw http.ResponseWriter, r *http.Request) {
		var page int
		if p := r.URL.Query().Get(pageParam); p != "" {
			if i, err := strconv.Atoi(p); err == nil {
				page = i
			}
		}

		var status int
		if p := r.URL.Query().Get(statusParam); p != "" {
			if i, err := strconv.Atoi(p); err == nil {
				status = i
			}
		}

		res, err := a.MessageList(&app.MessageListRequest{Page: page, Status: status})
		if err != nil {
			renderError(rw, err, http.StatusInternalServerError)
			return
		}

		data := left.IndexData{
			Messages: res.Messages,
			Paginate: paginate.New(*r.URL, pageParam, res.PageMin, res.Page, res.PageMax),
		}
		render(rw, w, left.IndexPage, &data)
	}
}
