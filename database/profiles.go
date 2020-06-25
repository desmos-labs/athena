package database

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/profile"
	dbtypes "github.com/desmos-labs/djuno/database/types"
	"github.com/rs/zerolog/log"
)

// SaveUserIfNotExisting creates a new user having the given address if it does not exist yet.
// Upon creating the user, returns that.
// If any error is raised during the process, returns that.
func (db DesmosDb) SaveUserIfNotExisting(address sdk.AccAddress) (*dbtypes.ProfileRow, error) {
	// Insert the user
	sqlStmt := `INSERT INTO profile (address) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := db.sqlx.Exec(sqlStmt, address.String())
	if err != nil {
		return nil, err
	}

	return db.GetUserByAddress(address)
}

// SaveProfile saves the given profile into the database, replacing any existing info.
// Returns the inserted row or an error if something goes wrong.
func (db DesmosDb) SaveProfile(profile profile.Profile) (*dbtypes.ProfileRow, error) {
	log.Info().
		Str("module", "profiles").
		Str("moniker", profile.Moniker).
		Str("creator", profile.Creator.String()).
		Msg("saving profile")

	sqlStmt := `INSERT INTO profile (address, moniker, name, surname, bio, profile_pic, cover_pic) 
				VALUES ($1, $2, $3, $4, $5, $6, $7) 
				ON CONFLICT (address) DO UPDATE 
				    SET address = excluded.address, 
				        moniker = excluded.moniker, 
				        name = excluded.name, 
				        surname = excluded.surname,
				        bio = excluded.bio,
				        profile_pic = excluded.profile_pic,
				        cover_pic = excluded.cover_pic`

	var profilePic, coverPic *string
	if profile.Pictures != nil {
		profilePic = profile.Pictures.Profile
		coverPic = profile.Pictures.Cover
	}

	_, err := db.Sql.Exec(
		sqlStmt,
		profile.Creator.String(), profile.Moniker, profile.Name, profile.Surname, profile.Bio,
		profilePic, coverPic,
	)
	if err != nil {
		return nil, err
	}

	return db.GetUserByAddress(profile.Creator)
}

// DeleteProfile allows to delete the profile of the user having the given address
func (db DesmosDb) DeleteProfile(address sdk.AccAddress) error {
	updatedProfile := profile.NewProfile(address)
	_, err := db.SaveProfile(updatedProfile)
	return err
}

// ______________________________________

// ExecuteQueryAndGetFirstUserRow executes the given query with the specified arguments
// and returns the first matched row.
func (db DesmosDb) ExecuteQueryAndGetFirstUserRow(query string, args ...interface{}) (*dbtypes.ProfileRow, error) {
	var rows []dbtypes.ProfileRow
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

// GetUserById returns the user having the specified id. If not found returns nil instead.
func (db DesmosDb) GetUserById(id *uint64) (*dbtypes.ProfileRow, error) {
	return db.ExecuteQueryAndGetFirstUserRow(`SELECT * FROM profile WHERE id = $1`, id)
}

// GetUserByAddress returns the user row having the given address.
// If the user does not exist yet, returns nil instead.
func (db DesmosDb) GetUserByAddress(address sdk.AccAddress) (*dbtypes.ProfileRow, error) {
	return db.ExecuteQueryAndGetFirstUserRow(`SELECT * FROM profile WHERE address = $1`, address.String())
}
