package serve

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

var METHOD_ALL    = []string{"GET", "POST", "PUT", "DELETE"}
var METHOD_GET    = []string{"GET"}
var METHOD_POST   = []string{"POST"}
var METHOD_PUT    = []string{"PUT"}
var METHOD_DELETE = []string{"DELETE"}

/*
 * Router requirements:
 * 1. Handle basic, text only routes. E.g.: /index
 * 2. Handle routes based on method: GET, POST, PUT, DELETE. Also multiple methods should be allowed: {GET, POST}
 * 3. Handle placeholder values: /wish/:id/
 */

func (s *Serve) ServeRoute(w http.ResponseWriter, r *http.Request) {
	//Not logged in urls
	switch {
	case match("/login", METHOD_ALL, r):
		s.loginHandler(w, r)
		return
	case match("/register", METHOD_GET, r):
		s.registerHandlerGet(w, r)
		return
	case match("/register", METHOD_POST, r):
		s.registerHandlerPost(w, r)
		return
	}
	//Ensure we are logged in for the following urls
	loggedIn, user := s.getLoggedInUserOrRedirect(w, r)
	if !loggedIn {
		return
	}

	var id uint64
	var userName string
	switch {
	case match("/logout", METHOD_ALL, r):
		s.logoutHandler(w, r)
	case match("/", METHOD_ALL, r):
		s.indexHandler(user, w, r)
	case match("/wish", METHOD_GET, r):
		td := templateData{"", "", "", "", false, false, ""}
		s.renderEditWishTemplate(w, td)
	case match("/wish", METHOD_POST, r):
		s.newWishPostHandler(user, w, r)
	case match("/wish/:id/delete", METHOD_GET, r, &id):
		s.deleteWishHandler(user, w, r, id)
	case match("/wish/:id/edit", METHOD_GET, r, &id):
		s.editWishGetHandler(user, w, r, id)
	case match("/wish/:id/edit", METHOD_POST, r, &id):
		s.editWishPostHandler(user, w, r, id)
	case match("/list/:user", METHOD_ALL, r, &userName):
		s.otherUserHandler(user, w, r, userName)
	case match("/list/:user/wish/:id/reserve", METHOD_ALL, r, &userName, &id):
		s.reserveWishHandler(user, w, r, userName, id)
	case match("/list/:user/wish/:id/unreserve", METHOD_ALL, r, &userName, &id):
		s.unreserveWishHandler(user, w, r, userName, id)
	default:
		s.notFoundHandler(w, r)
	}
}

func match(expectedPattern string, expectedMethods []string, r *http.Request, vars ...any) bool {
	methodOk := false
	for _, method := range expectedMethods {
		if method == r.Method {
			methodOk = true
		}
	}
	if !methodOk {
		return false
	}
	path := r.URL.Path
	pathSlices := strings.Split(path, "/")
	patternSlices := strings.Split(expectedPattern, "/")

	if len(pathSlices) != len(patternSlices) {
		return false
	}

	argumentIdx := 0

	for i, patternPart := range patternSlices {
		if len(patternPart) > 0 && patternPart[0] == ':' {
			//We have a variable to fill
			if argumentIdx >= len(vars) {
				log.Print("Not enough arguments provided for URL pattern")
				return false
			}
			switch p := vars[argumentIdx].(type) {
			case *string:
				*p = pathSlices[i]
			case *int:
				n, err := strconv.Atoi(pathSlices[i])
				if err != nil {
					return false
				}
				*p = n
			case *uint64:
				n, err := strconv.ParseUint(pathSlices[i], 10, 64)
				if err != nil {
					return false
				}
				*p = n
			case *int64:
				n, err := strconv.ParseInt(pathSlices[i], 10, 64)
				if err != nil {
					return false
				}
				*p = n
			default:
				log.Println("vars must be {*string, *int, *uint64, *int64}")
				return false
			}
			argumentIdx += 1
		} else {
			if patternPart != pathSlices[i] {
				return false
			}
		}
	}
	if argumentIdx != len(vars) {
		//Not all arguments were filled
		return false
	}
	return true
}
