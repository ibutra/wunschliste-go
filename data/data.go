package data

import (
  bolt "go.etcd.io/bbolt"
)

const (
  BUCKETNAME = "wishlist"
)

type Data struct {
  db *bolt.DB
}
