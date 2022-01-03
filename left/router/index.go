package router

import (
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/ItsNotGoodName/smtpbridge/left/web"
	"github.com/ItsNotGoodName/smtpbridge/pkg/paginate"
)

func handleIndexGet(t *web.Templater, a *app.App) http.HandlerFunc {
	type Data struct {
		Messages []app.Message
		Paginate paginate.Paginate
	}

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

		pag := paginate.New(*r.URL, param, res.PageMin, res.Page, res.PageMax)
		data := Data{Messages: res.Messages, Paginate: pag}

		t.Render(web.PageIndex, rw, data)
	}
}
