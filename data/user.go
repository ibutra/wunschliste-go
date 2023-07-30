package data

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
	"golang.org/x/crypto/argon2"
)

const (
	SALT_SIZE       = 30
	ITERATION_COUNT = 1
	MEMORY          = 512 * 1024 // 512MB
	CPU_COUNT       = 4
	HASH_LENGTH     = 128
)

var userBucketName []byte = []byte("user")

var UserNotExistingError = errors.New("User does not exist")

//Values must not be changed! Only public for saving to database
type User struct {
	Name string
	Hash []byte
	Salt []byte
	data *Data
}

func (d *Data) CreateUser(name string, password string) (*User, error) {
	user, err := d.GetUser(name)
	if err != nil && err != UserNotExistingError {
		return nil, err
	}
	if user != nil {
		return nil, errors.New("User already present")
	}
	salt := make([]byte, SALT_SIZE)
	_, err = rand.Read(salt)
	if err != nil {
		log.Printf("Failed to generate salt: %v", err)
		return nil, err
	}
	hash := hashPassword(password, salt)
	user = &User{
		Name: name,
		Hash: hash,
		Salt: salt,
		data: d,
	}
	err = user.save()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) save() error {
	err := u.data.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(userBucketName)
		if err != nil {
			return err
		}
		json, err := json.Marshal(u)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(u.Name), json)
	})
	return err
}

func (u *User) ChangePassword(password string) error {
  salt := make([]byte, SALT_SIZE)
  _, err := rand.Read(salt)
  if err != nil {
    log.Printf("Failed to generate salt: %v", err)
    return err
  }
  hash := hashPassword(password, salt)
  u.Salt = salt
  u.Hash = hash
  return u.save()
}

func (u *User) CheckPassword(password string) bool {
	enteredHash := hashPassword(password, u.Salt)
	return subtle.ConstantTimeCompare(enteredHash, u.Hash) == 1
}

func (u *User) String() string {
	return fmt.Sprintf("User: %v", u.Name)
}

func (d *Data) GetUser(name string) (*User, error) {
	var user *User = nil
	err := d.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userBucketName))
		if bucket == nil {
			return UserNotExistingError
		}
		jsonData := bucket.Get([]byte(name))
		if jsonData == nil {
			return UserNotExistingError
		}
		user = &User{
			data: d,
		}
		return json.Unmarshal(jsonData, user)
	})
	return user, err
}

func hashPassword(password string, salt []byte) []byte {
	hash := argon2.IDKey([]byte(password), salt, ITERATION_COUNT, MEMORY, CPU_COUNT, HASH_LENGTH)
	return hash
}
