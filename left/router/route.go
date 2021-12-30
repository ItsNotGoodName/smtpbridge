package router

func (s *Router) route() {
	s.r.Get(s.attachmentURI+"*", s.GetAttachments())
	s.r.Get("/", s.GetIndex())
}
