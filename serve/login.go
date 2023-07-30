package serve

import (
	"net/http"
)

func loginHandler(s *Serve, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.templates.ExecuteTemplate(w, "login", nil)
		return
	}
	name := r.PostFormValue("name")
	password := r.PostFormValue("password")

	if name != "" && password != "" {
		user, err := s.data.GetUser(name)
		if err == nil {
			if user.CheckPassword(password) {
				//Create session
				w.Write([]byte("Login successful"))
				return
			}
		}
	}
		s.templates.ExecuteTemplate(w, "login", nil)
}
