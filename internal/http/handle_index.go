package http

import "net/http"

func (s Service) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, _ := s.t.ParseFiles(pageIndex)
		t.Execute(w, nil)
	}
}
