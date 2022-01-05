package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Database contains the connection pool of SQL.
// Server should call Initialize before usage.
// Thanks to Go maintainers, the concurrency support is inbuilt
// so we do not need to manage connections on our own.
type Database struct {
	pool *sql.DB
}

// Initialize will open a database at given dataSourceName.
func (db *Database) Initialize(dataSourceName string) error {

	pool, err := sql.Open("postgres", dataSourceName)
	if err == nil {
		db.pool = pool
	}
	return err
}

// Close is meant to free up resources, to be called when the
// server wants to shut down.
func (db *Database) Close() error {
	return db.pool.Close()
}
