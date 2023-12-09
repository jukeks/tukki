package db_test

import (
	"testing"

	"github.com/jukeks/tukki/db"
	"github.com/jukeks/tukki/lib/testhelpers"
	"github.com/thanhpk/randstr"
)

func TestDB(t *testing.T) {
	dbDir := testhelpers.EnsureTempDirectory("test-tukki-" + randstr.String(10))
	database := db.NewDatabase(dbDir)

	key := randstr.String(10)
	value := randstr.String(16 * 1024)
	database.Put(key, value)

	storedValue, err := database.Get(key)
	if err != nil {
		t.Fatalf("failed to get value: %v", err)
	}

	if storedValue != value {
		t.Fatalf("stored value does not match: %s != %s", storedValue, value)
	}

	database.Close()

	// reopen database
	database = db.NewDatabase(dbDir)
	storedValue, err = database.Get(key)
	if err != nil {
		t.Fatalf("failed to get value: %v", err)
	}

	if storedValue != value {
		t.Fatalf("stored value does not match: %s != %s", storedValue, value)
	}

	database.Close()
}
