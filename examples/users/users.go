package users

import (
	"context"

	"github.com/romeovs/preql/sql"
)

// A User in the database.
type User struct {
	// Username is the user's username.
	Username string `sql:"username"`

	// Email is the user's email address, normalized.
	Email string `sql:"email"`
}

// A Repository that holds users.
type Repository struct {
	sql sql.Interface
}

// FindUserByUsername finds a user with the given username.
// Returns nil if no user is found.
//
// sql:
//   SELECT * FROM users
//     WHERE lower(username) = lower(:username)
//     LIMIT 1
func (r Repository) FindUserByUsername(ctx context.Context, username string) (*User, error) {
	return toUser(preqlFindUserByUsername(r.sql, ctx, username))
}
