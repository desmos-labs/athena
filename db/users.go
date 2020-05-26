package db

import (
	"database/sql"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/profile"
)

// UserRow represents a single PostgreSQL row containing the data of a user
type UserRow struct {
	Id         *uint64        `db:"id"`
	Address    string         `db:"address"`
	Moniker    sql.NullString `db:"moniker"`
	Name       sql.NullString `db:"name"`
	Surname    sql.NullString `db:"surname"`
	Bio        sql.NullString `db:"bio"`
	ProfilePic sql.NullString `db:"profile_pic"`
	CoverPic   sql.NullString `db:"cover_pic"`
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
	sqlStmt := `INSERT INTO "user" (address) VALUES ($1)`
	_, err = db.sqlx.Exec(sqlStmt, address.String())
	if err != nil {
		return nil, err
	}

	return db.GetUserByAddress(address)
}

// UpsertProfile saves the given profile into the database, replacing any existing info.
// Returns the inserted row or an error if something goes wrong.
func (db DesmosDb) UpsertProfile(account profile.Profile) (*UserRow, error) {
	row, err := db.GetUserByAddress(account.Creator)
	if err != nil {
		return nil, err
	}

	if row != nil {
		// If a row already exists, and some fields are null, we are going to replace them
		// with the existing values to avoid any unwanted overwrite

		if account.Name == nil && row.Name.Valid {
			account.Name = &row.Name.String
		}

		if account.Surname == nil && row.Surname.Valid {
			account.Surname = &row.Surname.String
		}

		if account.Bio == nil && row.Bio.Valid {
			account.Bio = &row.Bio.String
		}

		pictures := account.Pictures
		if pictures == nil {
			pictures = &profile.Pictures{Profile: nil, Cover: nil}
		}

		if pictures.Profile == nil && row.ProfilePic.Valid {
			pictures.Profile = &row.ProfilePic.String
		}

		if pictures.Cover == nil && row.CoverPic.Valid {
			pictures.Cover = &row.CoverPic.String
		}
	}

	return db.saveProfile(account)
}

// saveProfile saves the given profile by replacing any existing data about it, or creating
// a new entry if it does not exist yet.
func (db DesmosDb) saveProfile(profile profile.Profile) (*UserRow, error) {
	row, err := db.GetUserByAddress(profile.Creator)
	if err != nil {
		return nil, err
	}

	// Create if not exists
	if row == nil {
		sqlStmt := `INSERT INTO "user" (address, moniker, name, surname, bio, profile_pic, cover_pic) 
					VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err := db.Sql.Exec(
			sqlStmt,
			profile.Creator.String(), profile.Moniker, profile.Name, profile.Surname, profile.Bio,
			profile.Pictures.Profile, profile.Pictures.Cover,
		)
		if err != nil {
			return nil, err
		}

		return db.GetUserByAddress(profile.Creator)
	}

	// Update
	sqlStmt := `UPDATE "user" 
				SET moniker = $1, name = $2, surname = $3, bio = $4, profile_pic = $5, cover_pic = $6
				WHERE address = $7`
	_, err = db.Sql.Exec(sqlStmt,
		profile.Moniker, profile.Name, profile.Surname, profile.Bio, profile.Pictures.Profile, profile.Pictures.Cover,
		profile.Creator.String())
	if err != nil {
		return nil, err
	}

	return db.GetUserByAddress(profile.Creator)
}

// DeleteProfile allows to delete the profile of the user having the given address
func (db DesmosDb) DeleteProfile(address sdk.AccAddress) error {
	updatedProfile := profile.NewProfile(address)
	_, err := db.UpsertProfile(updatedProfile)
	return err
}
