package db

import "github.com/jukeks/tukki/internal/db"

type Database = db.Database
type Config = db.Config
type Pair = db.Pair
type Cursor = db.Cursor

var ErrKeyNotFound = db.ErrKeyNotFound

func GetDefaultConfig() Config {
	return db.GetDefaultConfig()
}

func OpenDatabase(dbDir string, config Config) (*Database, error) {
	return db.OpenDatabaseWithConfig(dbDir, config)
}
