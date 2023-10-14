package main

import (
	"flag"
	"log"

	"github.com/ibutra/wunschliste-go/data"
	"github.com/ibutra/wunschliste-go/serve"
)

var dataFilePath = flag.String("file", "db.bolt", "Database file to use")

func main() {
	flag.Parse()

	d, err := data.OpenData(*dataFilePath)

	if err != nil {
		log.Panic("Failed to open Database: ", err)
	}
	defer d.Close()

	serve, err := serve.NewServe(&d)
	if err != nil {
		log.Panic("Failed to serve: ", err)
	}

	err = serve.Serve()
	if err != nil {
		log.Panic("Failed to serve ", err)
	}
}
