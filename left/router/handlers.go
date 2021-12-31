package router

import (
	"embed"
	"html/template"
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/smtpbridge/app"
	"github.com/ItsNotGoodName/smtpbridge/dto"
)

func (s *Router) GetAttachments() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		http.StripPrefix(s.attachmentURI, http.FileServer(http.FS(s.a.AttachmentGetFS()))).ServeHTTP(rw, r)
	}
}

//go:embed template
var templateFS embed.FS

type Paginate struct {
	Page    int
	PageMax int
}

func (p Paginate) HasPrev() bool {
	return p.Page > 1
}

func (p Paginate) FirstLink() string {
	return "?page=" + strconv.Itoa(1)
}

func (p Paginate) LastLink() string {
	return "?page=" + strconv.Itoa(p.PageMax)
}

func (p Paginate) PrevLink() string {
	return "?page=" + strconv.Itoa(p.Page-1)
}

func (p Paginate) HasNext() bool {
	return p.Page < p.PageMax
}

func (p Paginate) NextLink() string {
	return "?page=" + strconv.Itoa(p.Page+1)
}

func (s *Router) GetIndex() http.HandlerFunc {
	type Data struct {
		Messages []dto.Message
		Paginate Paginate
	}

	index := template.Must(template.ParseFS(templateFS, "template/index.html"))

	return func(rw http.ResponseWriter, r *http.Request) {
		var page int
		if p := r.URL.Query().Get("page"); p != "" {
			if i, err := strconv.Atoi(p); err == nil {
				page = i
			}
		}

		res, err := s.a.MessageList(&app.MessageListRequest{AttachmentPath: s.attachmentURI, Page: page})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = index.Execute(rw, Data{Messages: res.Messages, Paginate: Paginate{Page: res.Page, PageMax: res.PageMax}})
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
