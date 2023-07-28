package serve

import (
	"fmt"
	"net/http"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		templates.ExecuteTemplate(w, "login", nil)
		return
	}
	name := r.PostFormValue("name")
	password := r.PostFormValue("password")

	if name == "" || password == "" {
		templates.ExecuteTemplate(w, "login", nil)
		return
	}

	w.Write([]byte(fmt.Sprintf("Name: %v PW: %v", name, password)))
}
