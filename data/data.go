package data

import (
	bolt "go.etcd.io/bbolt"
)

const (
	DATABASE_FILE = "db.bolt"
)

type Data struct {
	db *bolt.DB
}

func OpenData() (Data, error) {
	db, err := bolt.Open(DATABASE_FILE, 0666, nil)
	if err != nil {
		return Data{}, err
	}
	data := Data{
		db: db,
	}
	return data, nil
}

func (d *Data) Close() {
	d.db.Close()
}
