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

	_, err = d.CreateUser("Lukas", "blub")
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
	if user.CheckPassword("abcde") {
		log.Printf("Password check OK")
	} else {
		log.Printf("Password check failed")
	}

	wish, err := user.CreateWish("XxXxX", 123.98, "blub")
	if err != nil {
		log.Panic("Failed to create wish: ", err)
	}

	log.Printf("%v\n", wish)

	wishs, err := user.GetWishs()
	if err != nil {
		log.Panic("Failed to get wishs: ", err)
	}
	log.Printf("%v\n", wishs)

	err = wish.Delete()
	if err != nil {
		log.Panic("Failed to delete wish: ", err)
	}

	wishs, err = user.GetWishs()
	if err != nil {
		log.Panic("Failed to get wishs: ", err)
	}
	log.Printf("%v\n", wishs)

	err = wishs[0].Reserve(user)
	if err != nil {
		log.Panic("Failed to reserve wish: ", err)
	}
	log.Printf("%v\n", wishs)
}
