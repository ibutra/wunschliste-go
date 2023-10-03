package serve

import (
	"fmt"
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
		serve.renderUserList(loggedInUser, w, r, true, false, loggedInUser.Name)
	} else {
		user, err := serve.data.GetUser(userName)
		if err != nil {
			log.Println("User dosn't exist")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		serve.renderUserList(user, w, r, false, true, loggedInUser.Name)
	}
}

func (serve *Serve) renderUserList(user data.User, w http.ResponseWriter, r *http.Request, canEdit bool, canReserve bool, lookingUser string) {
	var ti templateInfo
	ti.CanEdit = canEdit
	ti.CanReserve = canReserve
	allWishs, err := user.GetWishs()
	if err != nil {
		log.Println("Failed to get wishs for user: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	wishs := make([]data.Wish, 0)
	//Only show wishes that are reserved by us or not reserved
	if lookingUser == user.Name {
		wishs = allWishs
	} else {
		for _, wish := range(allWishs) {
			if wish.Reserved == "" || wish.Reserved == lookingUser {
				wishs = append(wishs, wish)
			}
		}
	}
	ti.Wishs = wishs
	serve.renderNavbar(w)
	if err := serve.templates.ExecuteTemplate(w, "wishlist", ti); err != nil {
		log.Println(err)
	}
}

func (serve *Serve) reserveWishHandler(loggedInUser data.User, w http.ResponseWriter, r *http.Request, userName string, wishId uint64) {
	user, err := serve.data.GetUser(userName)
	if err != nil {
		log.Println("Failed to get user: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	wish, err := user.GetWishWithId(wishId)
	if err != nil {
		log.Println("Failed to get wish for user: ", err)
		http.Redirect(w, r, fmt.Sprintf("/list/%v", userName), http.StatusTemporaryRedirect)
		return
	}
	if wish.Reserved == "" {
		err = wish.Reserve(&loggedInUser)
		if err != nil {
			log.Println("Failed to reserve wish: ", err)
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/list/%v", userName), http.StatusTemporaryRedirect)
}

func (serve *Serve) unreserveWishHandler(loggedInUser data.User, w http.ResponseWriter, r *http.Request, userName string, wishId uint64) {
	user, err := serve.data.GetUser(userName)
	if err != nil {
		log.Println("Failed to get user: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	wish, err := user.GetWishWithId(wishId)
	if err != nil {
		log.Println("Failed to get wish for user: ", err)
		http.Redirect(w, r, fmt.Sprintf("/list/%v", userName), http.StatusTemporaryRedirect)
		return
	}
	err = wish.Unreserve()
	if err != nil {
		log.Println("Failed to unreserve wish: ", err)
	}
	http.Redirect(w, r, fmt.Sprintf("/list/%v", userName), http.StatusTemporaryRedirect)
}
