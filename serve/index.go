package serve

import (
	"net/http"
	"fmt"
)

func indexHandler(serve *Serve, w http.ResponseWriter, r *http.Request) {
	loggedIn, user := serve.getLoggedInUserOrRedirect(w, r)
	if !loggedIn {
		return
	}
	serve.templates.ExecuteTemplate(w, "index", fmt.Sprintf("Hallo %v! Willkommen bei Wunschliste", user.Name))
}
