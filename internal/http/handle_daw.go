package http

import "net/http"

func (s Service) handleDaw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, _ := s.t.ParseFiles(pageIndex, pageDaw)
		t.Execute(w, nil)
	}
}
