package serve

import (
	"net/http"
	"encoding/base64"
	"log"
)

func (s *Serve) logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		log.Print(err)
		return
	}
	secret, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		log.Print(err)
		return
	}
	session, err := s.data.GetSessionFromSecret(secret)
	if err != nil {
		log.Print(err)
		return
	}
	err = session.Delete()
	if err != nil {
		log.Print(err)
		return
	}
}
