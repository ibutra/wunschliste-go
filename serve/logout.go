package serve

import (
	"net/http"
	"encoding/base64"
	"log"
)

func logoutHandler(s *Serve, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
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
