package serve

import (
	"log"
	"net/http"
)

func (s *Serve) registerHandlerGet(w http.ResponseWriter, r *http.Request) {
	err := s.templates.ExecuteTemplate(w, "register", nil)
	if err != nil {
		log.Println(err)
	}
}

func (s *Serve) registerHandlerPost(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("username")
	password := r.PostFormValue("password")
	_, err := s.data.CreateUser(name, password)
	if err != nil {
		log.Println(err)
		s.templates.ExecuteTemplate(w, "register", "Failure registering")
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
