package data

import (
  "golang.org/x/crypto/argon2"
  "crypto/rand"
  "crypto/subtle"
  "log"
  "errors"
  bolt "go.etcd.io/bbolt"
)

const (
  SALT_SIZE = 30
  ITERATION_COUNT = 1
  MEMORY = 512 * 1024 // 512MB
  CPU_COUNT = 4 
  HASH_LENGTH = 128

  BUCKETNAME = "user"
)

type User struct {
  Name string
  hash []byte
  salt []byte
}

func (d Data) CreateUser(name string, password string) (User, error) {
  user := User{}
  err := d.db.Update(func (tx *bolt.Tx) error {
    //Check if user present
    bucket,err := tx.CreateBucketIfNotExists([]byte(BUCKETNAME));
    if err != nil {
      return err
    }
    userPresent := bucket.Get([]byte(name)) != nil

    if userPresent {
      return errors.New("User already present")
    }

    salt := make([]byte, SALT_SIZE)
    _, err = rand.Read(salt)
    if err != nil {
      log.Printf("Failed to generate salt: %v", err)
      return err
    }
    hash := hashPassword(password, salt)
    user = User{
      Name: name,
      hash: hash,
      salt: salt,
    }
    //TODO: save to database
    return nil
  })
  return user, err
}

func (u User) CheckPassword(password string) bool {
  enteredHash := hashPassword(password, u.salt)
  return subtle.ConstantTimeCompare(enteredHash, u.hash) == 1
}

func (d Data) GetUser(name string) (User, error) {
  user := User{}
  err := d.db.View(func (tx *bolt.Tx) error {
    bucket := tx.Bucket([]byte(BUCKETNAME))
    if bucket == nil {
      return errors.New("Bucket not present")
    }
    
    return nil
  })
  return user, err
}

func hashPassword(password string, salt []byte) []byte {
  hash := argon2.IDKey([]byte(password), salt, ITERATION_COUNT, MEMORY, CPU_COUNT, HASH_LENGTH)
  return hash
}
