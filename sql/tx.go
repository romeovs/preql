package sql

import (
	"context"
	"database/sql"
)

// Tx implements Interface by wrapping a database/sql.Tx.
type Tx sql.Tx

// ExecContext implements Interface.
func (tx *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error) {
	return (*sql.Tx)(tx).ExecContext(ctx, query, args...)
}

// QueryContext implements Interface.
func (tx *Tx) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	return (*sql.Tx)(tx).QueryContext(ctx, query, args...)
}

// Begin implements Interface.
// No new transaction is opened in the driver, since a transaction is already open.
func (tx *Tx) Begin(ctx context.Context, opts *TxOptions) (Interface, error) {
	return tx, nil
}
