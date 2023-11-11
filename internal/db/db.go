package db

import (
	"path/filepath"

	"github.com/enkodr/machina/internal/path"
	"go.etcd.io/bbolt"
	bolt "go.etcd.io/bbolt"
)

// Represents a structure for the DB connection
type DB struct {
	db     *bolt.DB
	bucket string
}

// Creates a new DB connection
func NewDB(bucket string) (*DB, error) {
	db := &DB{
		bucket: bucket,
	}
	var err error
	dbPath := filepath.Join(path.GetPath(path.DatabaseFile))

	// Open the database file
	db.db, err = bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	// Initialise the database and create the bucket if it doesn't exist
	db.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return nil
	})

	return db, nil
}

func (db *DB) Put(key string, value []byte) error {
	// Store the data
	err := db.db.Update(func(tx *bbolt.Tx) error {
		// Create a bucket.
		bkt := tx.Bucket([]byte(db.bucket))

		return bkt.Put([]byte(key), value)
	})

	return err
}

func (db *DB) Get(key string) ([]byte, error) {
	val := []byte{}
	db.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(db.bucket))
		val = bkt.Get([]byte(key))
		return nil
	})
	return val, nil
}

func (db *DB) Close() {
	db.db.Close()
}
