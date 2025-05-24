package db

import "github.com/jukeks/tukki/internal/db"

func GetDefaultConfig() db.Config {
	return db.GetDefaultConfig()
}

func OpenDatabase(dbDir string, config db.Config) (*db.Database, error) {
	return db.OpenDatabaseWithConfig(dbDir, config)
}
