package db

import "github.com/jukeks/tukki/internal/db"

type Database = db.Database
type Config = db.Config

func GetDefaultConfig() Config {
	return db.GetDefaultConfig()
}

func OpenDatabase(dbDir string, config Config) (*Database, error) {
	return db.OpenDatabaseWithConfig(dbDir, config)
}
