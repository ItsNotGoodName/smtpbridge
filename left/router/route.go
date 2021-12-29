package router

import "net/http"

func (s *Router) route() {
	s.r.Get("/attachments/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/attachments/", http.FileServer(http.Dir(s.AttachmentsPath))).ServeHTTP(w, r)
	})
}
