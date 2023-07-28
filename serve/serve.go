package serve

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed templates/*.html
var templatesFS embed.FS

var templates *template.Template

func Serve() error {
	t, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return err
	}
	templates = t

	mux := http.NewServeMux()
	
	mux.HandleFunc("/login", loginHandler)
	err = http.ListenAndServe(":8080", mux)
	return err
}
