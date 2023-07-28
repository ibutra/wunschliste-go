package main

import (
	"log"

	"github.com/ibutra/wunschliste-go/data"
	"github.com/ibutra/wunschliste-go/serve"
)

func main() {
	d, err := data.OpenData()

	if err != nil {
		log.Panic("Failed to open Database: ", err)
	}
	defer d.Close()

	err = serve.Serve()
	if err != nil {
		log.Panic("Failed to serve: ", err)
	}

}
