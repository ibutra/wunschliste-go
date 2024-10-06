package data

import (
	"encoding/json"
	"errors"

	bolt "go.etcd.io/bbolt"
)

var settingsBucketName []byte = []byte("Settings")
var settingsKey []byte = []byte("Settings")
var SettingsNotExistingError = errors.New("Settings not saved")

type Settings struct {
  RegisterClosed bool
}

func (d *Data) SetSettings(settings Settings) error {
  err := d.db.Update(func(tx *bolt.Tx) error {
    bucket, err := tx.CreateBucketIfNotExists(settingsBucketName)
    if err != nil {
      return err
    }
    json, err := json.Marshal(settings)
    if err != nil {
      return err
    }
    return bucket.Put(settingsKey, json)
  })
  return err
}

func (d *Data) GetSettings() (Settings, error) {
  var settings Settings
  err := d.db.View(func(tx *bolt.Tx) error {
    bucket := tx.Bucket(settingsBucketName)
    if bucket == nil {
      return SettingsNotExistingError
    }
    jsonData := bucket.Get(settingsKey)
    if jsonData == nil {
      return SettingsNotExistingError
    }
    return json.Unmarshal(jsonData, &settings)
  })
  return settings, err
}
