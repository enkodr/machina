package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) (*DB, func()) {
	t.Helper()

	// Create a temporary file for the database
	tmpfile, err := os.CreateTemp("/tmp/", "testdb-*.db")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}

	// Initialize DB
	db, err := NewDB("testbucket")
	if err != nil {
		t.Fatalf("could not create db: %v", err)
	}

	// Cleanup function to close and remove the database file
	cleanup := func() {
		db.Close()
		os.Remove(tmpfile.Name())
	}

	return db, cleanup
}

func TestNewDB(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	assert.NotNil(t, db)
}

func TestPutAndGet(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	err := db.Put("testkey", []byte("testvalue"))
	assert.Nil(t, err)

	val, err := db.Get("testkey")
	assert.Nil(t, err)
	assert.Equal(t, "testvalue", string(val))
}
