package serve

import (
	"log"
	"net/http"
)

func (s *Serve) registerHandlerGet(w http.ResponseWriter, r *http.Request) {
  settings, err := s.data.GetSettings()
  if err != nil {
    log.Println(err)
  }
  if settings.RegisterClosed {
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
    return
  }
	err = s.templates.ExecuteTemplate(w, "register", nil)
	if err != nil {
		log.Println(err)
	}
}

func (s *Serve) registerHandlerPost(w http.ResponseWriter, r *http.Request) {
  settings, err := s.data.GetSettings()
  if err != nil {
    log.Println(err)
  }
  if settings.RegisterClosed {
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
    return
  }
	name := r.PostFormValue("username")
	password := r.PostFormValue("password")
	user, err := s.data.CreateUser(name, password)
	if err != nil {
		log.Println(err)
		s.templates.ExecuteTemplate(w, "register", "Fehler beim Registrieren")
		return
	}
	//If this is the first user, approve automatically and make admin
	if s.data.GetUserCount() == 1 {
		if err = user.Approve(); err != nil {
			log.Println(err)
		}
		if err = user.SetAdmin(true); err != nil {
			log.Println(err)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
