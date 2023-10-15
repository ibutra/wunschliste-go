package serve

import (
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/ibutra/wunschliste-go/data"
)

var sessionTimeout time.Duration = 100 * 365 * 24 * 60 * 60 * 1000 * 1000 * 1000 //100 years
var sessionCookieName = "wunschliste-session"

func (s *Serve) loginHandler(w http.ResponseWriter, r *http.Request) {
	renderLoginTemplate(s, w, "")
	return
}

func (s *Serve) loginHandlerPost(w http.ResponseWriter, r *http.Request) {
	if loggedIn, _ := s.getLoggedInUser(r); loggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		log.Println("User alreaddy logged in tried to access login page")
		return
	}
	name := r.PostFormValue("username")
	password := r.PostFormValue("password")

	if name == "" || password == "" {
		renderLoginTemplate(s, w, "Es muss ein Name und Passwort angegeben werden.")
		return
	}

	user, err := s.data.GetUser(name)
	if err != nil {
		log.Print(err)
		renderLoginTemplate(s, w, "Falsche Login-Daten.")
		return
	}
	if !user.Approved {
		renderLoginTemplate(s, w, "Sie wurden bisher nicht freigeschaltet.")
		return
	}
	if !user.CheckPassword(password) {
		renderLoginTemplate(s, w, "Falsche Login-Daten.")
		return
	}
	err = login(user, w)
	if err != nil {
		log.Print(err)
		renderLoginTemplate(s, w, "Fehler beim Einloggen.")
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func renderLoginTemplate(s *Serve, w http.ResponseWriter, message string) {
	if message != "" {
		s.templates.ExecuteTemplate(w, "login", struct{ Message string }{Message: message})
	} else {
		s.templates.ExecuteTemplate(w, "login", nil)
	}
}

func login(user data.User, w http.ResponseWriter) error {
	//Create session
	session, err := user.CreateSession(sessionTimeout)
	if err != nil {
		return err
	}
	cookie := http.Cookie{
		Name:     sessionCookieName,
		Value:    base64.StdEncoding.EncodeToString(session.Secret),
		Expires:  session.ValidUntil,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	}
	http.SetCookie(w, &cookie)
	return nil
}

func (s *Serve) getLoggedInUser(r *http.Request) (bool, data.User) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		if err != http.ErrNoCookie {
			log.Print(err)
		}
		return false, data.User{}
	}
	secret, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		log.Print(err)
		return false, data.User{}
	}
	session, err := s.data.GetSessionFromSecret(secret)
	if err != nil {
		if err != data.NoActiveSessionError {
			log.Print(err)
		}
		return false, data.User{}
	}
	user, err := s.data.GetUser(session.User)
	if err != nil {
		log.Print(err)
		return false, data.User{}
	}
	return true, user
}

func (s *Serve) getLoggedInUserOrRedirect(w http.ResponseWriter, r *http.Request) (bool, data.User) {
	loggedIn, user := s.getLoggedInUser(r)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	return loggedIn, user
}
