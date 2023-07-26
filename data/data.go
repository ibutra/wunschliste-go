package data

import (
	"strconv"
	"strings"

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

func convertUInt64ToByteArray(i uint64) []byte {
	return []byte(strconv.FormatUint(i, 10))
}

func convertByteArrayToUint64(b []byte) (uint64, error) {
	return strconv.ParseUint(string(b), 10, 64)
}

func (d *Data) String() string {
	var sb strings.Builder
  d.db.View(func (tx *bolt.Tx) error {
		tx.ForEach(func (name []byte, bucket *bolt.Bucket) error {
			BucketString(name, bucket, 0, &sb)
			return nil
		})
		return nil
	})
	return sb.String()
}

func BucketString(name []byte, bucket *bolt.Bucket, depth int, sb *strings.Builder) {
	for i := 0; i<depth; i++ {
		sb.WriteString("  ")
	}
	sb.WriteString("Bucket: ")
	sb.Write(name)
	bucket.ForEach(func (k []byte, v []byte) error {
		if v == nil {
			return nil
		}
		sb.WriteString("\n")
		for i := 0; i < depth + 1; i++ {
			sb.WriteString("  ")
		}
		sb.Write(k)
		sb.WriteString(": ")
		sb.Write(v)
		return nil
	})
	bucket.ForEachBucket(func (k []byte) error {
		subBucket := bucket.Bucket(k)
		if subBucket != nil {
			sb.WriteString("\n")
			BucketString(k, subBucket, depth + 1, sb)
		}
		return nil
	})
	sb.WriteString("\n")
}
