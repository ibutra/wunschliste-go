package serve

import (
	"net/http"
	"log"
)

func indexHandler(serve *Serve, w http.ResponseWriter, r *http.Request) {
	loggedIn, user := serve.getLoggedInUserOrRedirect(w, r)
	if !loggedIn {
		return
	}
	wishs, err := user.GetWishs()
	if err != nil {
		log.Println("Failed to get wishs: ", err)
	}
	serve.templates.ExecuteTemplate(w, "index", wishs)
}
