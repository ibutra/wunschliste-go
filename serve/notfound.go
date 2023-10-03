package serve

import (
	"net/http"

	"github.com/ibutra/wunschliste-go/data"
)

func (s *Serve) notFoundHandler(loggedInUser data.User, w http.ResponseWriter, r *http.Request) {
	s.renderNavbar(loggedInUser, w)
	s.templates.ExecuteTemplate(w, "404", nil)
}
