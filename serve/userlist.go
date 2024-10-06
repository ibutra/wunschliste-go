package serve

import (
	"log"
	"net/http"

	"github.com/ibutra/wunschliste-go/data"
)

type userlistTemplateData struct {
	Message string
	Users []data.User
	LoggedInUser data.User
  RegisterClosed bool
}

func (s *Serve) serveUserList(loggedInUser data.User, w http.ResponseWriter, r *http.Request) {
	if !loggedInUser.Admin {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	var ti userlistTemplateData
	users, err := s.data.GetUsers()
	if err != nil {
		log.Println(err)
		ti.Message = "Failed to get all users"
	}
	ti.Users = users
	ti.LoggedInUser = loggedInUser
  settings, err := s.data.GetSettings()
  if err != nil {
    log.Println(err)
  }
  ti.RegisterClosed = settings.RegisterClosed
	s.renderNavbar(loggedInUser, w);
	err = s.templates.ExecuteTemplate(w, "userlist", ti)
	if err != nil {
		log.Println(err)
	}
}

func (s *Serve) approveUserHandler(loggedInUser data.User, w http.ResponseWriter, r *http.Request, userName string) {
	if !loggedInUser.Admin {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
	user, err := s.data.GetUser(userName)
	if err != nil {
		log.Println(err)
		return
	}
	err = user.Approve()
	if err != nil {
		log.Println(err)
		return
	}
}

func (s *Serve) deleteUserHandler(loggedInUser data.User, w http.ResponseWriter, r *http.Request, userName string) {
	if !loggedInUser.Admin {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
	user, err := s.data.GetUser(userName)
	if err != nil {
		log.Println(err)
		return
	}
	err = user.Delete()
	if err != nil {
		log.Println(err)
		return
	}
}

func (s *Serve) openRegistrationHandler(loggedInUser data.User, w http.ResponseWriter, r *http.Request) {
  if !loggedInUser.Admin {
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
  }
	http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
  settings := data.Settings{RegisterClosed: false};
  s.data.SetSettings(settings)
}

func (s *Serve) closeRegistrationHandler(loggedInUser data.User, w http.ResponseWriter, r *http.Request) {
  if !loggedInUser.Admin {
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
  }
	http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
  settings := data.Settings{RegisterClosed: true};
  s.data.SetSettings(settings)
}
