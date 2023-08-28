package serve

import (
	"net/http"

	"github.com/ibutra/wunschliste-go/data"
)

/*
 * Router requirements:
 * 1. Handle basic, text only routes. E.g.: /index
 * 2. Handle routes based on method: GET, POST, PUT, DELETE. Also multiple methods should be allowed: {GET, POST}
 * 3. Handle placeholder values: /wish/:id/
 * 4. More specific routes should have precedence
 */

type HandleFunc func (*Serve, data.User, http.ResponseWriter, *http.Request, ...interface{})

type Route struct {
	Method string
	Pattern string
	Handler HandleFunc
}

var routes = []Route{
	// {"GET", "/newWish", newWishHandler},
}

func ServeRoute(w http.ResponseWriter, r *http.Request) {

}

func match(path string, method string, route Route) bool {

	return false
}
