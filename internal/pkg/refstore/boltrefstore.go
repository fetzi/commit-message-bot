package refstore

import (
	bolt "go.etcd.io/bbolt"
)

// BoltRefStore struct
type BoltRefStore struct {
	db *bolt.DB
}

// NewBoltRefStore creates a new RefStore instance with a database handle
func NewBoltRefStore(dbPath string) (*BoltRefStore, error) {
	db, err := bolt.Open(dbPath, 0600, nil)

	if err != nil {
		return nil, err
	}

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("refstore"))

		if b == nil {
			tx.CreateBucket([]byte("refstore"))
		}

		return nil
	})

	return &BoltRefStore{db}, nil
}

// Put stores the given key value pair in the database
func (r *BoltRefStore) Put(key, value string) {
	r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("refstore"))
		return b.Put([]byte(key), []byte(value))
	})
}

// Exists checks if the given key exists in the database
func (r *BoltRefStore) Exists(key string) bool {
	var record string

	r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("refstore"))
		v := b.Get([]byte(key))

		if v != nil {
			record = string(v)
		}

		return nil
	})

	return record != ""
}
