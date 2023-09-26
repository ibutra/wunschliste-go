package serve

import (
	"log"
	"net/http"

	"github.com/ibutra/wunschliste-go/data"
)

type navbarData struct {
	Users []data.User
	Msg   string
}

func (s *Serve) renderNavbar(w http.ResponseWriter) {
	var nd  navbarData
	users, err := s.data.GetUsers()
	if err != nil {
		log.Println(err)
		nd.Msg = "Error with data backend"
	} else {
		nd.Users = users
	}
	s.templates.ExecuteTemplate(w, "navbar", nd)
}
