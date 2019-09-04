package sql

import (
	"context"
	"database/sql"
)

// DB implements Interface by wrapping database/sql.DB.
type DB sql.DB

// ExecContext implements Interface.
func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error) {
	return (*sql.DB)(db).ExecContext(ctx, query, args...)
}

// QueryContext implements Interface.
func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	return (*sql.DB)(db).QueryContext(ctx, query, args...)
}

// Begin implements Interface.
func (db *DB) Begin(ctx context.Context, opts *TxOptions) (Interface, error) {
	tx, err := (*sql.DB)(db).BeginTx(ctx, opts)
	return (*Tx)(tx), err
}

// Exec executes the query in the database.
func (db *DB) Exec(query string, args ...interface{}) (Result, error) {
	return db.ExecContext(context.Background(), query, args...)
}
