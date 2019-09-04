// Package sql package provides a uniform interface for sql
// operations on sql databases.
//
// Most notably Interface unifies operations in the top-level database (sql.DB) and
// operations in a transaction (sql.Tx).
package sql

import (
	"database/sql"
)

type (
	// TxOptions holds the transaction options to be used in DB.Begin.
	TxOptions = sql.TxOptions

	// Rows is the result of a query. Its cursor starts before the first row
	// of the result set. Use Next to advance from row to row.
	Rows = sql.Rows

	// A Result summarizes an executed SQL command.
	Result = sql.Result
)

var (
	// ErrNoRows is the error that gets returned if a query results in no rows.
	ErrNoRows = sql.ErrNoRows
)

// Open opens a database specified by its database driver name and a
// driver-specific data source name, usually consisting of at least a
// database name and connection information.
//
// Most users will open a database via a driver-specific connection
// helper function that returns a *DB. No database drivers are included
// in the Go standard library. See https://golang.org/s/sqldrivers for
// a list of third-party drivers.
//
// Open may just validate its arguments without creating a connection
// to the database. To verify that the data source name is valid, call
// Ping.
//
// The returned DB is safe for concurrent use by multiple goroutines
// and maintains its own pool of idle connections. Thus, the Open
// function should be called just once. It is rarely necessary to
// close a DB.
func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return (*DB)(db), nil
}
