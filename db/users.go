package db

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PostRow represents a single PostgreSQL row containing the data of a Post
type UserRow struct {
	Id      *uint64 `db:"id"`
	Address string  `db:"address"`
}

// GetUserById returns the user having the specified id. If not found returns nil instead.
func (db DesmosDb) GetUserById(id *uint64) (*UserRow, error) {
	return db.ExecuteQueryAndGetFirstUserRow(`SELECT * FROM "user" WHERE id = $1`, id)
}

// GetUserByAddress returns the user row having the given address.
// If the user does not exist yet, returns nil instead.
func (db DesmosDb) GetUserByAddress(address sdk.AccAddress) (*UserRow, error) {
	return db.ExecuteQueryAndGetFirstUserRow(`SELECT * FROM "user" WHERE address = $1`, address.String())
}

// ExecuteQueryAndGetFirstUserRow executes the given query with the specified arguments
// and returns the first matched row.
func (db DesmosDb) ExecuteQueryAndGetFirstUserRow(query string, args ...interface{}) (*UserRow, error) {
	var rows []UserRow
	err := db.sqlx.Select(&rows, query, args...)
	if err != nil {
		return nil, err
	}

	// No users found, return nil
	if len(rows) == 0 {
		return nil, nil
	}

	return &rows[0], nil
}

// SaveUserIfNotExisting creates a new user having the given address if it does not exist yet.
// Upon creating the user, returns that.
// If any error is raised during the process, returns that.
func (db DesmosDb) SaveUserIfNotExisting(address sdk.AccAddress) (*UserRow, error) {
	user, err := db.GetUserByAddress(address)
	if err != nil {
		return nil, err
	}

	// User already existing, do nothing
	if user != nil {
		return user, nil
	}

	// Insert the user
	var id *uint64
	sqlStmt := `INSERT INTO "user" (address) VALUES ($1) RETURNING id`
	err = db.Sql.QueryRow(sqlStmt, address.String()).Scan(&id)
	if err != nil {
		return nil, err
	}

	row := UserRow{Id: id, Address: address.String()}
	return &row, nil
}
