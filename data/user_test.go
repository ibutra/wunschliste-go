package data

import (
  "testing"
  "crypto/rand"
)

func TestCreateUser(t *testing.T) {
  name := "User1"
  password := make([]byte, 90)
  _, err := rand.Read(password)
  if err != nil {
    t.Fatal(err)
  }
  user, err := CreateUser(name, string(password))
  if err != nil {
    t.Fatal(err)
  }
  if user.Name != name {
    t.Error("Username not equal:", user.Name, name)
  }
  
  if !user.CheckPassword(string(password)) {
    t.Error("Failed to auth password")
  }
}
