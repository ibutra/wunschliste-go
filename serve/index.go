package serve

import (
	"log"
	"net/http"

	"github.com/ibutra/wunschliste-go/data"
)

type templateInfo struct {
	Wishs []data.Wish
	CanEdit bool
	CanReserve bool
}

func indexHandler(serve *Serve, user data.User, w http.ResponseWriter, r *http.Request) {
	wishs, err := user.GetWishs()
	if err != nil {
		log.Println("Failed to get wishs: ", err)
	}
	ti := templateInfo {
		Wishs: wishs,
		CanEdit: true,
		CanReserve: false,
	}
	if err := serve.templates.ExecuteTemplate(w, "index", ti); err != nil {
		log.Println(err)
	}
}
