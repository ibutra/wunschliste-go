package main

import (
	"fmt"
	"log"

	"github.com/ibutra/wunschliste-go/data"
)

func main() {
  fmt.Println("Hello Worldy")
  d, err := data.OpenData()

  if err != nil {
    log.Panic("Failed to open Database")
    return
  }
  defer d.Close()

  user, err := d.CreateUser("Stefan", "abcde")
  if err != nil {
    log.Panic("Failed to create User")
    return
  }

  fmt.Print(user)

}
