package router

func (s *Router) route() {
	s.r.Get("/attachments/*", s.GetAttachments("/attachments/"))
	s.r.Get("/", s.GetIndex())
}
