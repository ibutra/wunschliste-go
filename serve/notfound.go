package serve

import (
	"net/http"
)

func (s *Serve) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	s.renderNavbar(w)
	s.templates.ExecuteTemplate(w, "404", nil)
}
