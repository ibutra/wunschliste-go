package serve

import (
	"log"
	"net/http"

	"github.com/ibutra/wunschliste-go/data"
)

type templateInfo struct {
	Wishs      []data.Wish
	CanEdit    bool
	CanReserve bool
}

func (serve *Serve) indexHandler(user data.User, w http.ResponseWriter, r *http.Request) {
	wishs, err := user.GetWishs()
	if err != nil {
		log.Println("Failed to get wishs: ", err)
	}
	ti := templateInfo{
		Wishs:      wishs,
		CanEdit:    true,
		CanReserve: false,
	}
	serve.renderNavbar(w)
	if err := serve.templates.ExecuteTemplate(w, "wishlist", ti); err != nil {
		log.Println(err)
	}
}

func (serve *Serve) otherUserHandler(loggedInUser data.User, w http.ResponseWriter, r *http.Request, userName string) {
	if loggedInUser.Name == userName {
		serve.renderUserList(loggedInUser, w, r, true, false)
	} else {
		user, err := serve.data.GetUser(userName)
		if err != nil {
			log.Println("User dosn't exist")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		serve.renderUserList(user, w, r, false, true)
	}
}

func (serve *Serve) renderUserList(user data.User, w http.ResponseWriter, r *http.Request, canEdit bool, canReserve bool) {
	var ti templateInfo
	ti.CanEdit = canEdit
	ti.CanReserve = canReserve
	wishs, err := user.GetWishs()
	if err != nil {
		log.Println("Failed to get wishs for user: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	ti.Wishs = wishs
	serve.renderNavbar(w)
	if err := serve.templates.ExecuteTemplate(w, "wishlist", ti); err != nil {
		log.Println(err)
	}
}
