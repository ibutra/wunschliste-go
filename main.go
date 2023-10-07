package main

import (
	"flag"
	"fmt"
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

	if !handleArguments(&d) {
		return
	}

	if _, err := d.GetUser("Stefan"); err == data.UserNotExistingError {
		log.Println("Default user not present. Exiting")
		return
	}

	fmt.Print(d.String())

	serve, err := serve.NewServe(&d)
	if err != nil {
		log.Panic("Failed to serve: ", err)
	}

	err = serve.Serve()
	if err != nil {
		log.Panic("Failed to serve ", err)
	}

}

//Returns false if execution should not continue
func handleArguments(d *data.Data) bool{
	var cmd string = ""
	flag.StringVar(&cmd, "cmd", "", "Command for the executable: CreateUser, Approve, MakeAdmin")
	var user string = ""
	flag.StringVar(&user, "user", "", "User for the given command")
	var pw string = ""
	flag.StringVar(&pw, "password", "", "Password for the given command")
	
	flag.Parse()
	switch (cmd){
	case "CreateUser":
		if user == "" || pw == "" {
			log.Println("You must provide user and password")
			return false
		}
		user, err := d.CreateUser(user, pw)
		if err != nil {
			log.Println(err)
			return false
		}
		err = user.Approve()
		if err != nil {
			log.Println(err)
			return false
		}
		return false
	case "Approve":
		if user == "" {
			log.Println("You must provide a user")
			return false
		}
		user, err := d.GetUser(user)
		if err != nil {
			log.Println(err)
			return false
		}
		err = user.Approve()
		if err != nil {
			log.Println(err)
			return false
		}
		return false
	case "MakeAdmin":
		if user == "" {
			log.Println("You must provide a user")
			return false
		}
		user, err := d.GetUser(user)
		if err != nil {
			log.Println(err)
			return false
		}
		err = user.SetAdmin(true)
		if err != nil {
			log.Println(err)
			return false
		}
		return false
	}
	return true
}
