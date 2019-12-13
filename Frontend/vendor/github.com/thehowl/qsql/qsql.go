// Package qsql implements SQL queries for the lazy, in the good ol' hashtable
// or list of hashtable format.
package qsql

import (
	"database/sql"
	"errors"
)

// ErrDBIsNil is returned when the *sql.DB inside DB is nil.
var ErrDBIsNil = errors.New("qsql: db is nil")

// DB wraps an sql.DB around a custom DB.
//
// If you're hardcore and want to create one without New(),
// you should &qsql.DB{*db}.
type DB struct {
	sql.DB
}

// Exec behaves the same as sql.DB.Exec, however it does not wrap the last
// insert ID and rows affected into an interface.
func (d *DB) Exec(query string, params ...interface{}) (int, int, error) {
	res, err := d.DB.Exec(query, params...)
	lid, _ := res.LastInsertId()
	ra, _ := res.RowsAffected()
	return int(lid), int(ra), err
}

// ExecNoRes returns sql.DB.Exec without Result.
func (d *DB) ExecNoRes(query string, params ...interface{}) error {
	_, err := d.DB.Exec(query, params...)
	return err
}

// Query queries the database for multiple rows. See sql.DB.Query.
func (d *DB) Query(query string, params ...interface{}) ([]map[string]String, error) {
	return d.query(query, false, params...)
}

// QueryRow queries the database for one row. See sql.DB.QueryRow.
func (d *DB) QueryRow(query string, params ...interface{}) (map[string]String, error) {
	// sql.Row does not have .Columns(), so we can't really use db.QueryRow.
	// Instead, we use .query, telling it to return after the first row is extracted.
	m, err := d.query(query, true, params...)
	if len(m) > 0 {
		return m[0], err
	}
	return nil, err
}

func (d *DB) query(query string, only1 bool, params ...interface{}) ([]map[string]String, error) {
	rows, err := d.DB.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var returnSlice []map[string]String
	for rows.Next() {
		m, args := buildMapAndArgsSlice(cols)
		err := rows.Scan(args...)
		if err != nil {
			return nil, err
		}
		returnSlice = append(returnSlice, depointify(m))
		if only1 {
			return returnSlice, rows.Err()
		}
	}
	return returnSlice, rows.Err()
}

func buildMapAndArgsSlice(cols []string) (map[string]*string, []interface{}) {
	m := make(map[string]*string, len(cols))
	sl := make([]interface{}, 0, len(cols))
	for _, col := range cols {
		var newS string
		m[col] = &newS
		sl = append(sl, &newS)
	}
	return m, sl
}

func depointify(from map[string]*string) map[string]String {
	m := make(map[string]String, len(from))
	for k, v := range from {
		if v == nil {
			v = new(string)
		}
		m[k] = String(*v)
	}
	return m
}
