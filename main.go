package main

import (
	"log"
	"net/http"

	"github.com/ibutra/wunschliste-go/data"
)

func main() {
	d, err := data.OpenData()

	if err != nil {
		log.Panic("Failed to open Database: ", err)
	}
	defer d.Close()

	mux := http.NewServeMux()
	
	mux.HandleFunc("/", index)
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
}
