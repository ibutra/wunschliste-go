package data

import (
	"encoding/json"
	"strconv"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var wishBucketName []byte = []byte("wish")

type Wish struct {
	Name string
	Price float64
	Link string
	User string
	Reserved string
	data *Data
}

func (u *User) CreateWish(name string, price float64, link string) (*Wish, error) {
	var wish *Wish = nil
	err := u.data.db.Update(func (tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(wishBucketName)
		if err != nil {
			return err
		}
		bucket, err = bucket.CreateBucketIfNotExists([]byte(u.Name))
		if err != nil {
			return err
		}
		id, err := bucket.NextSequence()
		if err != nil {
			return err
		}
		wish = &Wish{
			Name: name,
			Price: price,
			Link: link,
			User: u.Name,
			data: u.data,
		}
		payload, err := json.Marshal(wish)
		if err != nil {
			return err
		}
		bucket.Put([]byte(strconv.FormatUint(id, 10)), payload)
		return nil
	})
	return wish, err
}

func (u *User) GetWishs() ([]*Wish, error) {
	wishs := make([]*Wish, 0)
	err := u.data.db.View(func (tx *bolt.Tx) error {
		bucket := tx.Bucket(wishBucketName)
		if bucket == nil {
			return nil
		}
		bucket = bucket.Bucket([]byte(u.Name))
		if bucket == nil {
			return nil
		}
		bucket.ForEach(func (k []byte, v []byte) error {
			var wish Wish
			json.Unmarshal(v, &wish)
			wish.data = u.data

			wishs = append(wishs, &wish)
			return nil
		})

		return nil
	})
	return wishs, err
}

func (w *Wish) String() string {
	return fmt.Sprintf("Wish: %v Price: %v Link: %v Owner: %v Reserved: %v", w.Name, w.Price, w.Link, w.User, w.Reserved)
}
