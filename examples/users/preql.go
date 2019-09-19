// Code generated by preql. DO NOT EDIT.

package users

import (
	"context"

	"github.com/romeovs/preql/sql"
)

// Scan scans the columns of the current row of rows into u.
// Expects rows.Next to be called at least once already.
func (u *User) Scan(rows *sql.Rows) error {
	names, err := rows.Columns()
	if err != nil {
		return err
	}

	// Undefined columns are ignored.
	var null interface{}

	scan := make([]interface{}, len(names))
	for i, name := range names {
		switch name {
		case "email":
			scan[i] = &u.Email
		case "username":
			scan[i] = &u.Username
		default:
			scan[i] = &null
		}
	}

	return rows.Scan(scan...)
}

// ScanOne scans one row from rows into u.
func (u *User) ScanOne(rows *sql.Rows) error {
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}
		return sql.ErrNoRows
	}

	return u.Scan(rows)
}

// toUser tries to scan the first row of rows into a new User.
// It checks err first to see if there were sql errors first and passes on the error.
// If err is sql.ErrNoRows, toUser returns nil, nil.
// This is a convenient helper to wrap preql query functions, eg.
//
//     u, err := toUser(preqlFindUser(r.sql, ctx))
//
func toUser(rows *sql.Rows, err error) (*User, error) {
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	u := new(User)
	err = u.Scan(rows)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// preqlFindUserByUsername implements the query associated with FindUserByUsername in it's doc comment.
func preqlFindUserByUsername(sql sql.Interface, ctx context.Context, username string) (*sql.Rows, error) {
	return sql.QueryContext(ctx, "SELECT * FROM users WHERE lower(username) = lower($1) LIMIT 1", username)
}
