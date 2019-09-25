# preql

[![](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/romeovs/preql)

A simple compile time query builder.

`preql` will search a package for sql queryies and generate simple helper
functions that perform that query at compile time so the runtime does not have
to parse queries anymore.

In addition `preql` can generate SQL scanners for the types in your code.

## Installation
```
go get -u github.com/romeovs/preql
```

## Usage

`preql` will not modify your code in any way, but can generate code that can
assist you in setting up sql queries and scanners.

### Scanners for a type

To generate scanners for a type, add `sql` struct tags to that type. For
example

```go
package user

type User struct {
  Username string `sql:"username"`
  Email    string `sql:"email"`
}
```

will generate the following functions in `preql.go`

```go
// Scan scans the columns of the current row of rows into u.
// Expects rows.Next to be called at least once already.
func (u *User) Scan (rows *sql.Rows) error {
  // Omitted
}

// ScanOne scans one row from rows into u.
func (u *User) ScanOne (rows *sql.Rows) error {
  // Omitted
}

// toUser tries to scan the first row of rows into a new User.
// It checks err first to see if there were sql errors first and passes on the error.
// If err is sql.ErrNoRows, toUser returns nil, nil.
// This is a convenient helper to wrap preql query functions, eg.
//
//     u, err := toUser(preqlFindUser(r.sql, ctx))
//
func toUser(rows *sql.Rows, err error) (*User, error) {
  // Omitted
}
```

These functions are written to be about as optimal as you can get. If you add
another field to the `User` struct definition, just run preql again and that
field will be scanned as well.

You can omit fields from being scanned by not setting the `sql` struct tag
on those fields.

### Bound queries

Say we have a `UserRepository` that wraps the `sql.DB` and will implement a
helper `FindUserByUsername` that can fetch a user from the database based on
theri username.

```go
package users

import (
  "database/sql"
)

type UserRepository struct {
  db *sql.DB
}

// FindUserByUsername finds a user by their username.
func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (*User, error) {
  // TODO: implement
  return nil, nil
}
```

You can hardcode the SQL query in a `r.db.QueryContext` call,
but this can become tedious because you need to remember to correctly number the
args and pass them in that order. This becomes really hard with big queries.

Instead we add a the SQL query to the comment like so:

```go
// FindUserByUsername finds a user by their username.
//
// sql:
//   SELECT * FROM users
//   WHERE lower(username) = lower(:username)
//   LIMIT 1
//
func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (*User, error) {
  // TODO: implement
  return nil, nil
}
```

If you run `preql` the user package now, `preql.go` will contain a function
`preqlFindUserByUsername`:

```go
// preqlFindUserByUsername implements the query associated with FindUserByUsername in it's doc comment.
func preqlFindUserByUsername(sql sql.Interface, ctx context.Context, username string) (*sql.Rows, error) {
  // Omitted
}
```

You can call this function and it will perform the provided query with the args
all in the correct order:

```go
func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (*User, error) {
  rows, err := preqlFindUserByUsername(r.db, ctx, username)
  // handle errors, scan user etc.
}
```

This makes writing sql queries a bit less tedious and using them a bit more
type-safe.

#### `context.Context`

`preql` detects wether or not the functions it generates need a
`context.Context` argument by looking at the orignal functions arguments.
Functions that don't have a `context.Context` argument will be using
`context.Background` internally.

#### Exec 
If you have a query that does not return results, use `exec:` in the comment
instread of `sql:`. For example:

```go
// CreateUser will create a new user in the system
//
// exec:
//   INSERT INTO users (
//     username,
//     email
//   ) VALUES (
//     lower(:username),
//     email
//   )
func (r *UserRepository) CreateUser(ctx context.Context, username string, email string) error {
  // Omitted
}
```

This will generate `preqlCreateUser` that uses exec under the hood.


#### Binding arguments

The sql query for a function can use all the arguments of the function
by their referencing name prefixed with `:`.

`context.Context` arguments are not passed to the query.

You can index structs by the name they have defined in the `sql` struct tag.
The `User` type above can be indexed by `.username` and `.email`, for example.

For example:
- `func(ctx context.Context, name string, limit uin64, offset uint64)` can use
  `:name`, `:limit`, and `:offset`
- `func(user *User)` can use `:user.username` and `:user.email`.


#### `sql.Interface`

The generated `preql*` functions do not accept a `database/sql.DB` but instead a
`github.com/romeovs/preql/sql.Interface` so they can be used in the global
database context as well as with transactions.

You will need to update the `UserRepository` as well:

```go
package users

import (
  "github.com/romeovs/preql/sql"
)

type UserRepository struct {
  db sql.Interface
}
```
