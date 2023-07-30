package data

import (
    bolt "go.etcd.io/bbolt"
    "crypto/rand"
    "time"
    "encoding/json"
    "errors"
)

const (
    SESSION_SECRET_SIZE = 128
)

var sessionBucketName []byte = []byte("session")

var NoActiveSessionError = errors.New("No active Session found")

type SessionSecret []byte

type Session struct {
    Secret SessionSecret
    ValidUntil time.Time
    User string
    data *Data
}

func (u *User) CreateSession(duration time.Duration) (Session, error) {
    err := u.data.garbageCollectSessions()
    if err != nil {
	return Session{}, err
    }
    secret := make(SessionSecret, SESSION_SECRET_SIZE)
    _, err = rand.Reader.Read(secret)
    if err != nil {
	return Session{}, err
    }
    var session Session
    err = u.data.db.Update(func (tx *bolt.Tx) error {
	bucket, err := tx.CreateBucketIfNotExists(sessionBucketName)
	if err != nil {
	    return err
	}
	session = Session {
	    Secret: secret,
	    ValidUntil: time.Now().Add(duration),
	    User: u.Name,
	    data: u.data,
	}
	jsonData, err := json.Marshal(session)
	if err != nil {
	    return err
	}
	return bucket.Put(secret, jsonData)
    })
    return session, err
}

func (d *Data) GetSessionFromSecret(secret SessionSecret) (Session, error) {
    err := d.garbageCollectSessions()
    if err != nil {
	return Session{}, err
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
    session.data = d
    //Don't neet to check the ValidUntil field as we garbage collected first
    return session, err
}

func (session *Session) Delete() error {
    err := session.data.garbageCollectSessions()
    if err != nil {
	return err
    }
    err = session.data.db.Update(func (tx *bolt.Tx) error {
	bucket := tx.Bucket(sessionBucketName)
	if bucket == nil {
	    return nil
	}
	return bucket.Delete(session.Secret)
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
