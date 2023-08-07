package main

import (
	"log"
	"fmt"

	"github.com/ibutra/wunschliste-go/data"
	"github.com/ibutra/wunschliste-go/serve"
)

func main() {
	d, err := data.OpenData()

	if err != nil {
		log.Panic("Failed to open Database: ", err)
	}
	defer d.Close()

	if _, err := d.GetUser("Stefan"); err == data.UserNotExistingError {
		user ,err := d.CreateUser("Stefan", "blub")
		if err != nil {
			log.Println("failed to create testuser", err)
			return
		}
		user.CreateWish("Test1", 12.4, "wunschliste.ibutra.com")
		user.CreateWish("Test2", 15.4, "wunschliste.ibutra.com")
	}

	fmt.Print(d.String())

	serve, err := serve.NewServe(&d)
	if err != nil {
		log.Panic("Failed to serve: ", err)
	}

	err = serve.Serve()

}
