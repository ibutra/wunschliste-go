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
	w.Write([]byte(fmt.Sprintf("Welcome %v", user.Name)))
}
