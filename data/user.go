package data

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"slices"

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
	Name      string
	Hash      []byte
	Salt      []byte
	Approved  bool
	Admin     bool
	VisibleTo []string //Name of other users we are visible to
	data      *Data
}

func (d *Data) CreateUser(name string, password string) (User, error) {
	user, err := d.GetUser(name)
	if err != nil && err != UserNotExistingError {
		return User{}, err
	}
	if err == nil {
		return User{}, errors.New("User already present")
	}
	salt := make([]byte, SALT_SIZE)
	_, err = rand.Read(salt)
	if err != nil {
		log.Printf("Failed to generate salt: %v", err)
		return User{}, err
	}
	otherUsers, err := d.GetUsers()
	if err != nil {
		log.Println("Failed to get other users: ", err)
		otherUsers = make([]User, 0)
	}
	otherUsersNames := make([]string, 0, len(otherUsers))
	for _, u := range(otherUsers) {
		otherUsersNames = append(otherUsersNames, u.Name)
	}
	hash := hashPassword(password, salt)
	user = User{
		Name: name,
		Hash: hash,
		Salt: salt,
		Approved: false,
		Admin: false,
		VisibleTo: otherUsersNames,
		data: d,
	}
	err = user.save()
	if err != nil {
		return User{}, err
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

func (u *User) SetAdmin(admin bool) error {
	u.Admin = admin
	return u.save()
}

func (u *User) Approve() error {
	u.Approved = true
	return u.save()
}

func (u *User) Delete() error {
	err := u.data.db.Update(func(tx *bolt.Tx) error {
		//Delete all wishs
		bucket := tx.Bucket(wishBucketName)
		if bucket != nil {
			if bucket.Bucket([]byte(u.Name)) != nil {
				err := bucket.DeleteBucket([]byte(u.Name))
				if err != nil {
					return err
				}
			}
		}
		//Delete user
		bucket = tx.Bucket([]byte(userBucketName))
		if bucket == nil {
			return UserNotExistingError
		}
		return bucket.Delete([]byte(u.Name))
	})
	return err
}

func (u *User) String() string {
	return fmt.Sprintf("User: %v", u.Name)
}

func (d *Data) GetUser(name string) (User, error) {
	var user User
	err := d.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userBucketName))
		if bucket == nil {
			return UserNotExistingError
		}
		jsonData := bucket.Get([]byte(name))
		if jsonData == nil {
			return UserNotExistingError
		}
		user = User{
			data: d,
		}
		return json.Unmarshal(jsonData, &user)
	})
	return user, err
}

func (d *Data) GetUserCount() int {
	count := 0
	err := d.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userBucketName))
		if bucket == nil {
			return nil
		}
		bucketStats := bucket.Stats()
		count = bucketStats.KeyN
		return nil
	})
	if err != nil {
		count = 0
		log.Println(err)
	}
	return count
}

func (d *Data) GetUsers() ([]User, error) {
	users := make([]User, 0)
	err := d.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userBucketName))
		if bucket == nil {
			return nil
		}
		err := bucket.ForEach(func (k, v []byte) error {
			user := User{
				data: d,
			}
			err := json.Unmarshal(v, &user)
			if err != nil {
				return err
			}
			users = append(users, user)
			return nil
		})
		return err
	})
	return users, err
}

func (u *User) GetVisibleUser() ([]User, error) {
	users := make([]User, 0)
	err := u.data.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(userBucketName))
		if bucket == nil {
			return UserNotExistingError
		}
		err := bucket.ForEach(func (k, v []byte) error {
			user := User{
				data: u.data,
			}
			err := json.Unmarshal(v, &user)
			if err != nil {
				return err
			}
			if slices.Contains(user.VisibleTo, u.Name) {
				users = append(users, user)
			}
			return nil
		})
		return err
	})
	return users, err
}

func hashPassword(password string, salt []byte) []byte {
	hash := argon2.IDKey([]byte(password), salt, ITERATION_COUNT, MEMORY, CPU_COUNT, HASH_LENGTH)
	return hash
}
