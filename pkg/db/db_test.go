package db

import "testing"

func TestOpen(t *testing.T) {
	t.Parallel()

	dbDir := t.TempDir()
	config := GetDefaultConfig()

	database, err := OpenDatabase(dbDir, config)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer database.Close()
}
