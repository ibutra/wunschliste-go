package data

import (
  "golang.org/x/crypto/argon2"
  "crypto/rand"
  "crypto/subtle"
  "log"
)

const (
  SALT_SIZE = 30
  ITERATION_COUNT = 1
  MEMORY = 512 * 1024 // 512MB
  CPU_COUNT = 4 
  HASH_LENGTH = 128
)

type User struct {
  id int64
  Name string
  hash []byte
  salt []byte
}

func CreateUser(name string, password string) (User, error) {
  salt := make([]byte, SALT_SIZE)
  _, err := rand.Read(salt)
  if err != nil {
    log.Printf("Failed to generate salt: %v", err)
    return User{}, err
  }
  hash := hashPassword(password, salt)
  user := User{
    id:  0, //TODO: implement
    Name: name,
    hash: hash,
    salt: salt,
  }
  return user, nil
}

func (u User) CheckPassword(password string) bool {
  enteredHash := hashPassword(password, u.salt)
  return subtle.ConstantTimeCompare(enteredHash, u.hash) == 1
}

func hashPassword(password string, salt []byte) []byte {
  hash := argon2.IDKey([]byte(password), salt, ITERATION_COUNT, MEMORY, CPU_COUNT, HASH_LENGTH)
  return hash
}
