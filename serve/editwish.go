package serve

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ibutra/wunschliste-go/data"
)

type templateData struct {
	Message    string
	Name       string
	Link       string
	PriceText  string
	NameRed    bool
	PriceRed   bool
	TargetLink string
}

type validatedInput struct {
	name  string
	link  string
	price float64
}

func (serve *Serve) newWishGetHandler(loggedInUser data.User, w http.ResponseWriter, r *http.Request) {
		td := templateData{"", "", "", "", false, false, ""}
		serve.renderEditWishTemplate(loggedInUser, w, td)
}

func (serve *Serve) newWishPostHandler(loggedInUser data.User, w http.ResponseWriter, r *http.Request) {
	if inputValid, input := serve.validateWishInput(loggedInUser, w, r, ""); inputValid {
		_, err := loggedInUser.CreateWish(input.name, input.price, input.link)
		if err != nil {
			log.Println(err)
			td := templateData{
				Message:    "Fehler beim Speichern des Wunsches. Administrator informiert",
				Name:       r.PostFormValue("name"),
				Link:       r.PostFormValue("link"),
				PriceText:  r.PostFormValue("price"),
				NameRed:    false,
				PriceRed:   false,
				TargetLink: "",
			}
			serve.renderEditWishTemplate(loggedInUser, w, td)
			return
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	//Don't need to render in this case as it is handled in the validateWishInput function
}

func (s *Serve) editWishGetHandler(loggedInUser data.User, w http.ResponseWriter, r *http.Request, wishId uint64) {
	wish, err := loggedInUser.GetWishWithId(wishId)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	td := templateData{
		Message:    "",
		Name:       wish.Name,
		Link:       wish.Link,
		PriceText:  strconv.FormatFloat(wish.Price, 'f', 2, 64),
		NameRed:    false,
		PriceRed:   false,
		TargetLink: fmt.Sprintf("/%v/edit", wishId),
	}
	s.renderEditWishTemplate(loggedInUser, w, td)
}

func (s *Serve) editWishPostHandler(loggedInUser data.User, w http.ResponseWriter, r *http.Request, wishId uint64) {
	target := fmt.Sprintf("/%v/edit", wishId)
	if inputValid, input := s.validateWishInput(loggedInUser, w, r, target); inputValid {
		td := templateData{
			Message:    "Fehler beim Speichern des Wunsches. Administrator informiert",
			Name:       r.PostFormValue("name"),
			Link:       r.PostFormValue("link"),
			PriceText:  r.PostFormValue("price"),
			NameRed:    false,
			PriceRed:   false,
			TargetLink: target,
		}
		wish, err := loggedInUser.GetWishWithId(wishId)
		if err != nil {
			log.Println("Edited wish not present for user", err)
			td.Message = "Fehler beim Speichern des Wunsches"
			s.renderEditWishTemplate(loggedInUser, w, td)
			return
		}
		wish.Name = input.name
		wish.Price = input.price
		wish.Link = input.link
		if err = wish.Save(); err != nil {
			log.Println(err)
			td.Message = "Fehler beim speichern des Wunsches"
			s.renderEditWishTemplate(loggedInUser, w, td)
			return
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	//Don't need to render in this case as it is handled in the validateWishInput function
}

func (s *Serve) validateWishInput(loggedInUser data.User, w http.ResponseWriter, r *http.Request, target string) (bool, validatedInput) {
	name := r.PostFormValue("name")
	link := r.PostFormValue("link")
	priceText := r.PostFormValue("price")
	td := templateData{
		Message:    "",
		Name:       name,
		Link:       link,
		PriceText:  priceText,
		NameRed:    false,
		PriceRed:   false,
		TargetLink: target,
	}
	if name == "" {
		td.Message = "Die Beschreibung darf nicht leer sein"
		td.NameRed = true
		s.renderEditWishTemplate(loggedInUser, w, td)
		return false, validatedInput{}
	}
	price, err := strconv.ParseFloat(priceText, 64)
	if err != nil {
		log.Println(err)
		td.Message = "Ungültiger Preis. Bitte nur Zahlen eingeben"
		td.PriceRed = true
		s.renderEditWishTemplate(loggedInUser, w, td)
		return false, validatedInput{}
	}
	if err != nil {
		log.Println(err)
		td.Message = "Fehler beim Speichern des Wunsches. Administrator informiert"
		s.renderEditWishTemplate(loggedInUser, w, td)
		return false, validatedInput{}
	}
	return true, validatedInput{
		name:  name,
		link:  link,
		price: price,
	}
}

func (s *Serve) deleteWishHandler(user data.User, w http.ResponseWriter, r *http.Request, wishId uint64) {
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

	wish, err := user.GetWishWithId(wishId)
	if err != nil {
		if err != data.WishNotPresent {
			log.Println(err)
		}
		return
	}
	if err = wish.Delete(); err != nil {
		log.Println(err)
		return
	}
}

func (s *Serve) renderEditWishTemplate(loggedInUser data.User, w http.ResponseWriter, td templateData) {
	s.renderNavbar(loggedInUser, w)
	if err := s.templates.ExecuteTemplate(w, "editWish", td); err != nil {
		log.Println(err)
	}
}
