package data

import (
    bolt "go.etcd.io/bbolt"
    "crypto/rand"
    "time"
    "encoding/json"
    "errors"
)

const (
    SESSION_SECRET_SIZE = 20
)

var sessionBucketName []byte = []byte("session")

var NoActiveSessionError = errors.New("No active Session found")

type SessionSecret []byte

type Session struct {
    ValidUntil time.Time
    User string
}

func (u *User) CreateSession(duration time.Duration) (SessionSecret, error) {
    err := u.data.garbageCollectSessions()
    if err != nil {
	return nil, err
    }
    secret := make(SessionSecret, SESSION_SECRET_SIZE)
    _, err = rand.Reader.Read(secret)
    if err != nil {
	return nil, err
    }
    err = u.data.db.Update(func (tx *bolt.Tx) error {
	bucket, err := tx.CreateBucketIfNotExists(sessionBucketName)
	if err != nil {
	    return err
	}
	session := Session {
	    ValidUntil: time.Now().Add(duration),
	    User: u.Name,
	}
	jsonData, err := json.Marshal(session)
	if err != nil {
	    return err
	}
	return bucket.Put(secret, jsonData)
    })
    return secret, err
}

func (d *Data) GetUserFromSession(secret SessionSecret) (User, error) {
    err := d.garbageCollectSessions()
    if err != nil {
	return User{}, err
    }
    var session Session
    err = d.db.View(func (tx *bolt.Tx) error {
	bucket := tx.Bucket(sessionBucketName)
	if bucket == nil {
	    return NoActiveSessionError
	}
	byteData := bucket.Get(secret)
	if byteData == nil {
	    return NoActiveSessionError
	}
	err := json.Unmarshal(byteData, &session)
	if err != nil {
	    return err
	}
	return nil
    })
    if err != nil {
	return User{}, err
    }
    //Don't neet to check the ValidUntil field as we garbage collect first
    return d.GetUser(session.User)
}

func (d *Data) DeleteSession(secret SessionSecret) error {
    err := d.garbageCollectSessions()
    if err != nil {
	return err
    }
    err = d.db.Update(func (tx *bolt.Tx) error {
	bucket := tx.Bucket(sessionBucketName)
	if bucket == nil {
	    return nil
	}
	return bucket.Delete(secret)
    })
    return err
}

func (d *Data) garbageCollectSessions() error {
    err := d.db.Update(func (tx *bolt.Tx) error {
	bucket := tx.Bucket(sessionBucketName)
	if bucket == nil {
	    return nil
	}
	//Go through all sessions and check if they are not expired
	oldSessions := make([]SessionSecret, 0, 10)
	var session Session
	err := bucket.ForEach(func (k []byte, v []byte) error {
	    err := json.Unmarshal(v, &session)
	    if err != nil {
		return err
	    }
	    if session.ValidUntil.Before(time.Now()) {
		oldSessions = append(oldSessions, k)
	    }
	    return nil
	})
	if err != nil {
	    return err
	}
	//Delete expired sessions
	for _, secret := range(oldSessions) {
	    err = bucket.Delete(secret)
	    if err != nil {
		return err
	    }
	}
	return nil
    })
    return err
}
