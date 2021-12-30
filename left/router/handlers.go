package router

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/ItsNotGoodName/smtpbridge/app"
)

func (s *Router) GetAttachments(prefix string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		http.StripPrefix(prefix, http.FileServer(http.FS(s.attachmentREPO.GetAttachmentFS()))).ServeHTTP(rw, r)
	}
}

//go:embed template
var templateFS embed.FS

func (s *Router) GetIndex() http.HandlerFunc {
	type Data struct {
		Messages []app.Message
	}

	index, err := template.ParseFS(templateFS, "template/index.html")
	if err != nil {
		log.Fatal("router.Router.GetIndex", err)
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		messages, err := s.messageSVC.List(10, 0)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		index.Execute(rw, Data{Messages: messages})
	}
}
