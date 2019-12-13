package qsql

import "database/sql"

// New creates a new *DB having an *sql.DB.
func New(db *sql.DB) *DB {
	if db == nil {
		return nil
	}
	return &DB{*db}
}

// Open behaves the same as sql.Open, but creates an *qsql.DB instead.
func Open(driverName, dsn string) (*DB, error) {
	db, err := sql.Open(driverName, dsn)
	return New(db), err
}
