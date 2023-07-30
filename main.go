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

	// d.CreateUser("Stefan", "blub")

	fmt.Print(d.String())

	serve, err := serve.NewServe(&d)
	if err != nil {
		log.Panic("Failed to serve: ", err)
	}

	err = serve.Serve()

}
