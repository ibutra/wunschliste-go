package serve

import (
	"net/http"
	"strconv"
	"strings"
	"log"
)

/*
 * Router requirements:
 * 1. Handle basic, text only routes. E.g.: /index
 * 2. Handle routes based on method: GET, POST, PUT, DELETE. Also multiple methods should be allowed: {GET, POST}
 * 3. Handle placeholder values: /wish/:id/
 * 4. More specific routes should have precedence
 */

func ServeRoute(w http.ResponseWriter, r *http.Request) {

}

func match(expectedPattern string, expectedMethod string, r *http.Request, vars ...any) bool {
	if expectedMethod != r.Method {
		return false
	}
	path := r.URL.Path
	pathSlices := strings.Split(path, "/")
	patternSlices := strings.Split(expectedPattern, "/")

	argumentIdx := 0

	for i, patternPart := range(patternSlices) {
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
			default:
				log.Println("vars must be *string or *int")
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
