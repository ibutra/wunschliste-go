package serve

import (
	"log"
	"net/http"
	"time"
	"encoding/base64"

	"github.com/ibutra/wunschliste-go/data"
)

var sessionTimeout time.Duration = 30 * 60 * 1000 * 1000 * 1000 //30 Minutes
var cookieName = "wunschliste-session"

func loginHandler(s *Serve, w http.ResponseWriter, r *http.Request) {
	if loggedIn, _ := s.getLoggedInUser(r); loggedIn {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if r.Method != http.MethodPost {
		renderLoginTemplate(s, w, "")
		return
	}
	name := r.PostFormValue("name")
	password := r.PostFormValue("password")

	if name == "" || password == "" {
		renderLoginTemplate(s, w, "You must provide name and password")
		return
	}

	user, err := s.data.GetUser(name)
	if err != nil {
		log.Print(err)
		renderLoginTemplate(s, w, "User not found")
		return
	}
	if !user.CheckPassword(password) {
		renderLoginTemplate(s, w, "Incorrect password")
		return
	}
	err = login(user, w)
	if err != nil {
		log.Print(err)
		renderLoginTemplate(s, w, "Error loggin in")
		return
	}
	log.Print("User logged in: ", user.Name)
	http.Redirect(w, r, "/", http.StatusFound)
}

func renderLoginTemplate(s *Serve, w http.ResponseWriter, message string) {
	if message != "" {
		s.templates.ExecuteTemplate(w, "login", struct {Message string}{Message: message})
	} else {
		s.templates.ExecuteTemplate(w, "login", nil)
	}
}

func login(user data.User, w http.ResponseWriter) error {
	//Create session
	secret, err := user.CreateSession(sessionTimeout)
	if err != nil {
		return err
	}
	cookie := http.Cookie {
		Name: cookieName,
		Value: base64.StdEncoding.EncodeToString(secret),
	}
	http.SetCookie(w, &cookie)
	return nil
}

func (s *Serve) getLoggedInUser(r *http.Request) (bool, data.User) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		log.Print(err)
		return false, data.User{}
	}
	secret, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		log.Print(err)
		return false, data.User{}
	}
	user, err := s.data.GetUserFromSession(secret)
	if err != nil {
		if err != data.NoActiveSessionError {
			log.Print(err)
		}
		return false, data.User{}
	}
	return true, user
}

func (s *Serve) getLoggedInUserOrRedirect(w http.ResponseWriter, r *http.Request) (bool, data.User) {
	loggedIn, user := s.getLoggedInUser(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	return loggedIn, user
}
