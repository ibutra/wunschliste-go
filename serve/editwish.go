package serve

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ibutra/wunschliste-go/data"
)

type templateData struct {
	Message string;
	Name string;
	Link string;
	PriceText string;
	NameRed bool;
	LinkRed bool;
}

func (serve *Serve) editWishPostHandler(user data.User, w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("name")
	link := r.PostFormValue("link")
	priceText := r.PostFormValue("price")
	if name == "" {
		td := templateData{
			Message: "Die Beschreibung darf nicht leer sein",
			Name: name,
			Link: link,
			PriceText: priceText,
			NameRed: true,
			LinkRed: false,
		}
		renderEditWishTemplate(serve, w, td)
		return
	}
	price, err := strconv.ParseFloat(priceText, 64)
	if err != nil {
		log.Println(err)
		td := templateData{
			Message: "Ungültiger Preis. Bitte nur Zahlen eingeben",
			Name: name,
			Link: link,
			PriceText: priceText,
			NameRed: false,
			LinkRed: true,
		}
		renderEditWishTemplate(serve, w, td)
		return
	}
	_, err = user.CreateWish(name, price, link)
	if err != nil {
		log.Println(err)
		renderEditWishTemplate(serve, w, templateData{"Fehler beim Speichern des Wunsches. Administrator informiert", name, link, priceText, false, false})
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func renderEditWishTemplate(serve *Serve, w http.ResponseWriter, td templateData) {
	if td.Message != "" {
		if err := serve.templates.ExecuteTemplate(w, "newWish", td); err != nil {
			log.Println(err)
		}
	} else {
		if err := serve.templates.ExecuteTemplate(w, "newWish", nil); err != nil {
			log.Println(err)
		}
	}
}
