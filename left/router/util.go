package router

import "net/http"

func multiplexAction(get, post, delete http.HandlerFunc) http.HandlerFunc {
	if get == nil {
		get = func(w http.ResponseWriter, r *http.Request) {}
	} else if post == nil {
		post = func(w http.ResponseWriter, r *http.Request) {}
	} else if delete == nil {
		delete = func(w http.ResponseWriter, r *http.Request) {}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		action := r.URL.Query().Get("action")
		if action == "delete" {
			delete(w, r)
		} else if action == "post" {
			post(w, r)
		} else {
			get(w, r)
		}
	}
}
