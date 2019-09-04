package sql

import (
	"context"
)

// Interface is the interface of SQL databases and transactions.
type Interface interface {
	// ExecContext executes a query without returning any rows.
	ExecContext(context.Context, string, ...interface{}) (Result, error)

	// QueryContext executes a query that returns rows, typically a SELECT.
	QueryContext(context.Context, string, ...interface{}) (*Rows, error)

	// Begin begins a new transaction and return the interface
	Begin(context.Context, *TxOptions) (Interface, error)
}
