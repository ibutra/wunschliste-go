package serve

import (
	"log"
	"net/http"

	"github.com/ibutra/wunschliste-go/data"
)

type navbarData struct {
	Users        []data.User
	Msg          string
	LoggedInUser data.User
}

func (s *Serve) renderNavbar(loggedInUser data.User, w http.ResponseWriter) {
	var nd  navbarData
	users, err := s.data.GetUsers()
	if err != nil {
		log.Println(err)
		nd.Msg = "Error with data backend"
	} else {
		nd.Users = users
	}
	nd.LoggedInUser = loggedInUser
	s.templates.ExecuteTemplate(w, "navbar", nd)
}
