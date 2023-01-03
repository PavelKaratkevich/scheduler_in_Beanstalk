package util

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
)

const (
	DefaultDatabaseName = "db.sqlite"
)

var (
	DataDir = getEnvDefault("DATA_DIR", "../data")
)

func ConnectDB() *sql.DB {
	db, err := sql.Open("sqlite3", filepath.Join(DataDir, DefaultDatabaseName))

	if err != nil {
		log.Fatal(err)
	}

	return db
}

// getEnvDefault return the value of the environment variable specified by name, or the defaultValue if not set
func getEnvDefault(name string, defaultValue string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}

	return defaultValue
}
