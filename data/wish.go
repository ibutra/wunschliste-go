package data

import (
	"encoding/json"
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var wishBucketName []byte = []byte("wish")

var WishBucketMissing = errors.New("Wish bucket doesn't exist")
var UserWishBucketMissing = errors.New("User's wish bucket doesn't exist")
var WishNotPresent = errors.New("No wish with this id present for given user")

type Wish struct {
	Name     string
	Price    float64
	Link     string
	User     string //owning user
	Reserved string //Userid who reserved
	Count		 int64 //How many are wished?
	Id       uint64 `json:"-"` //Must not be changed!
	data     *Data
}

func (u *User) CreateWish(name string, price float64, link string) (Wish, error) {
	var wish Wish
	err := u.data.db.Update(func(tx *bolt.Tx) error {
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
		wish = Wish{
			Name:  name,
			Price: price,
			Link:  link,
			User:  u.Name,
			Count: 1,
			data:  u.data,
			Id:    id,
		}
		payload, err := json.Marshal(wish)
		if err != nil {
			return err
		}
		bucket.Put(convertUInt64ToByteArray(id), payload)
		return nil
	})
	return wish, err
}

func (u *User) GetWishs() ([]Wish, error) {
	wishs := make([]Wish, 0)
	err := u.data.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(wishBucketName)
		if bucket == nil {
			return nil
		}
		bucket = bucket.Bucket([]byte(u.Name))
		if bucket == nil {
			return nil
		}
		bucket.ForEach(func(k []byte, v []byte) error {
			var wish Wish
			err := json.Unmarshal(v, &wish)
			if err != nil {
				return err
			}
			id, err := convertByteArrayToUint64(k)
			if err != nil {
				return err
			}
			wish.Id = id
			wish.data = u.data
			if wish.Count == 0 {
				wish.Count = 1;
			}

			wishs = append(wishs, wish)
			return nil
		})

		return nil
	})
	return wishs, err
}

func (u *User) GetWishWithId(id uint64) (Wish, error) {
	var wish Wish
	err := u.data.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(wishBucketName)
		if bucket == nil {
			return WishNotPresent
		}
		bucket = bucket.Bucket([]byte(u.Name))
		if bucket == nil {
			return WishNotPresent
		}
		key := convertUInt64ToByteArray(id)
		wishData := bucket.Get(key)
		if wishData == nil {
			return WishNotPresent
		}
		err := json.Unmarshal(wishData, &wish)
		if err != nil {
			return err
		}
		wish.data = u.data
		wish.Id = id
		return nil
	})
	return wish, err
}

func (w *Wish) Reserve(who *User) error {
	err := w.data.db.Update(func (tx *bolt.Tx) error {
	  bucket := tx.Bucket(wishBucketName)
		if bucket == nil {
			return WishBucketMissing
		}
		bucket = bucket.Bucket([]byte(w.User))
		if bucket == nil {
			return UserWishBucketMissing
		}
		oldReserved := w.Reserved
		w.Reserved = who.Name
		payload, err := json.Marshal(w)
		if err != nil {
			w.Reserved = oldReserved
			return err
		}
		return bucket.Put(convertUInt64ToByteArray(w.Id), payload)
	})
	return err
}

func (w *Wish) Delete() error {
	err := w.data.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(wishBucketName)
		if bucket == nil {
			return WishBucketMissing
		}
		bucket = bucket.Bucket([]byte(w.User))
		if bucket == nil {
			return UserWishBucketMissing
		}
		fmt.Println(w)
		return bucket.Delete(convertUInt64ToByteArray(w.Id))
	})
	return err
}

func (w *Wish) String() string {
	return fmt.Sprintf("Wish: %v Price: %v Link: %v Owner: %v Reserved: %v id: %v", w.Name, w.Price, w.Link, w.User, w.Reserved, w.Id)
}
