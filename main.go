package main

import (
	"log"

	"github.com/ibutra/wunschliste-go/data"
)

func main() {
	d, err := data.OpenData()

	if err != nil {
		log.Panic("Failed to open Database: ", err)
	}
	defer d.Close()

	_, err = d.CreateUser("Stefan", "abcde")
	if err != nil {
		log.Print("Failed to create User: ", err)
	} else {
		log.Print("Created user")
	}

	user, err := d.GetUser("Stefan")
	if err != nil {
		log.Panic("Failed to get User: ", err)
	}
	log.Printf("%v\n", user)

	wish, err := user.CreateWish("Brot", 0.12, "Bäcker")
	if err != nil {
		log.Panic("Failed to create wish: ", err)
	}

	log.Printf("%v\n", wish)

	wishs, err := user.GetWishs()
	if err != nil {
		log.Panic("Failed to get wishs: ", err)
	}
	log.Printf("%v\n", wishs)
}
