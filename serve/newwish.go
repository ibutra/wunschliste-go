package serve

import (
	"log"
	"strconv"
	"net/http"
)


func newWishHandler(serve *Serve, w http.ResponseWriter, r *http.Request) {
	loggedIn, user := serve.getLoggedInUserOrRedirect(w, r)
	if !loggedIn {
		return
	}
	if r.Method != http.MethodPost {
		renderNewWishTemplate(serve, w, "", "", "", "", false, false)
		return
	}
	name := r.PostFormValue("name")
	link := r.PostFormValue("link")
	priceText := r.PostFormValue("price")
	if name == "" {
		renderNewWishTemplate(serve, w, "Die Beschreibung darf nicht leer sein", name, link, priceText, true, false)
		return
	}
	price, err := strconv.ParseFloat(priceText, 64)
	if err != nil {
		log.Println(err)
		renderNewWishTemplate(serve, w, "Ungültiger Preis", name, link, priceText, false, true)
		return
	}
	_, err = user.CreateWish(name, price, link)
	if err != nil {
		log.Println(err)
		renderNewWishTemplate(serve, w, "Fehler beim Speichern des Wunsches. Administrator informiert", name, link, priceText, false, false)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func renderNewWishTemplate(serve *Serve, w http.ResponseWriter, message string, name string, link string, priceText string, nameInputRed bool, priceInputRed bool) {
	if message != "" {
		templateData := struct {
			Message string;
			Name string;
			Link string;
			PriceText string;
			NameRed bool;
			LinkRed bool;
		}{
			message,
			name,
			link,
			priceText,
			nameInputRed,
			priceInputRed,
		}
		if err := serve.templates.ExecuteTemplate(w, "newWish", templateData); err != nil {
			log.Println(err)
		}
	} else {
		if err := serve.templates.ExecuteTemplate(w, "newWish", nil); err != nil {
			log.Println(err)
		}
	}
}
