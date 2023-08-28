package serve

import (
	"net/http"
)

func (s *Serve) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	s.templates.ExecuteTemplate(w, "404", nil)
}
