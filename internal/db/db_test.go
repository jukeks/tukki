package db_test

import (
	"testing"

	"github.com/jukeks/tukki/internal/db"
	testutil "github.com/jukeks/tukki/tests/util"
	"github.com/thanhpk/randstr"
)

func TestDB(t *testing.T) {
	dbDir := testutil.EnsureTempDirectory("test-tukki-" + randstr.String(10))
	database := db.NewDatabase(dbDir)

	key := randstr.String(10)
	value := randstr.String(16 * 1024)
	database.Set(key, value)

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

func TestDelete(t *testing.T) {
	dbDir := testutil.EnsureTempDirectory("test-tukki-" + randstr.String(10))
	database := db.NewDatabase(dbDir)

	key := randstr.String(10)
	value := randstr.String(16 * 1024)
	database.Set(key, value)

	storedValue, err := database.Get(key)
	if err != nil {
		t.Fatalf("failed to get value: %v", err)
	}

	if storedValue != value {
		t.Fatalf("stored value does not match: %s != %s", storedValue, value)
	}

	database.Delete(key)

	_, err = database.Get(key)
	if err == nil {
		t.Fatalf("key should not exist anymore")
	}

	database.Close()
}
