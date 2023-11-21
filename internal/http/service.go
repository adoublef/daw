package http

import (
	"embed"
	"net/http"

	"github.com/adoublef/daw/internal/daw/lit"
	"github.com/adoublef/daw/template"
	"github.com/go-chi/chi/v5"
)

const (
	pageIndex = "index.html"
	pageDaw   = "daw.html"
)

//go:embed all:*.html
var fsys embed.FS
var FS = template.NewFS(fsys).Funcs(lit.FuncMap)

type Service struct {
	m *chi.Mux
	t *template.FS
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.ServeHTTP(w, r)
}

func NewService() *Service {
	s := &Service{
		m: chi.NewMux(),
		t: FS,
	}
	s.routes()
	return s
}

func (s Service) routes() {
	s.m.Get("/", s.handleIndex())
	s.m.Get("/daw", s.handleDaw())
	s.m.Handle("/daw/assets/*", lit.Handler)
}
