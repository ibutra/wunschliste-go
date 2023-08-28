package serve

import (
	"embed"
	"html/template"
	"net/http"
	"github.com/ibutra/wunschliste-go/data"
)

//go:embed templates/*.html
var templatesFS embed.FS

//go:embed static/*
var staticFS embed.FS

type Serve struct {
	templates *template.Template
	data *data.Data
	mux *http.ServeMux
}

type serveHandler func (*Serve, http.ResponseWriter, *http.Request)
type loggedInServeHandler func (*Serve, data.User, http.ResponseWriter, *http.Request)

func NewServe(data *data.Data) (*Serve, error) {
	t, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return &Serve{}, err
	}
	templates := t

	mux := http.NewServeMux()

	serve := &Serve {
		templates: templates,
		data: data,
		mux: mux,
	}

	//*******************
	// HANDLERS
	//*******************
	serve.addLoggedInHandler("/", indexHandler)
	serve.addHandler("/login", loginHandler)
	serve.addHandler("/logout", logoutHandler)
	serve.addLoggedInHandler("/newWish", newWishHandler)
	
	mux.Handle("/static/", http.FileServer(http.FS(staticFS)))

	return serve, nil
}

func (s *Serve) Serve() error {
	return http.ListenAndServe(":8080", s.mux)
}

func (s *Serve) addHandler(pattern string, handler serveHandler) {
	s.mux.HandleFunc(pattern, func (w http.ResponseWriter, r *http.Request) {
		handler(s, w, r)
	})
}

func (s *Serve) addLoggedInHandler(pattern string, handler loggedInServeHandler) {
	s.mux.HandleFunc(pattern, func (w http.ResponseWriter, r *http.Request) {
		loggedIn, user := s.getLoggedInUserOrRedirect(w, r)
		if !loggedIn {
			return
		}
		handler(s, user, w, r)
	})
}
