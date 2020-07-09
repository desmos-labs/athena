package database

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	profilestypes "github.com/desmos-labs/desmos/x/profiles"
	dbtypes "github.com/desmos-labs/djuno/database/types"
	"github.com/rs/zerolog/log"
)

// SaveUserIfNotExisting creates a new user having the given address if it does not exist yet.
// Upon creating the user, returns that.
// If any error is raised during the process, returns that.
func (db DesmosDb) SaveUserIfNotExisting(address sdk.AccAddress) error {
	// Insert the user
	sqlStmt := `INSERT INTO profile (address) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := db.Sqlx.Exec(sqlStmt, address.String())
	if err != nil {
		return err
	}

	_, err = db.GetUserByAddress(address)
	return err
}

// SaveProfile saves the given profilesTypes into the database, replacing any existing info.
// Returns the inserted row or an error if something goes wrong.
func (db DesmosDb) SaveProfile(profile profilestypes.Profile) error {
	log.Info().
		Str("module", "profiles").
		Str("dtag", profile.DTag).
		Str("creator", profile.Creator.String()).
		Msg("saving profile")

	sqlStmt := `INSERT INTO profile (address, dtag, moniker, bio, profile_pic, cover_pic, creation_date) 
				VALUES ($1, $2, $3, $4, $5, $6, $7) 
				ON CONFLICT (address) DO UPDATE 
				    SET address = excluded.address, 
				        dtag = excluded.dtag,
				        moniker = excluded.moniker, 
				        bio = excluded.bio,
				        profile_pic = excluded.profile_pic,
				        cover_pic = excluded.cover_pic,
				        creation_date = excluded.creation_date`

	var profilePic, coverPic *string
	if profile.Pictures != nil {
		profilePic = profile.Pictures.Profile
		coverPic = profile.Pictures.Cover
	}

	_, err := db.Sql.Exec(sqlStmt,
		profile.Creator.String(), profile.DTag, profile.Moniker, profile.Bio, profilePic, coverPic, profile.CreationDate)
	if err != nil {
		return err
	}

	_, err = db.GetUserByAddress(profile.Creator)
	return err
}

// DeleteProfile allows to delete the profilesTypes of the user having the given address
func (db DesmosDb) DeleteProfile(address sdk.AccAddress) error {
	sqlStmt := `UPDATE profile 
				SET dtag = $1, moniker = $2, bio = $3, profile_pic = $4, cover_pic = $5, creation_date = $6 
				WHERE address = $7`
	_, err := db.Sql.Exec(sqlStmt,
		nil, nil, nil, nil, nil, nil, address.String())
	return err
}

// ______________________________________

// executeQueryAndGetFirstUserRow executes the given query with the specified arguments
// and returns the first matched row.
func (db DesmosDb) executeQueryAndGetFirstUserRow(query string, args ...interface{}) (*dbtypes.ProfileRow, error) {
	var rows []dbtypes.ProfileRow
	err := db.Sqlx.Select(&rows, query, args...)
	if err != nil {
		return nil, err
	}

	// No users found, return nil
	if len(rows) == 0 {
		return nil, nil
	}

	return &rows[0], nil
}

// GetUserByAddress returns the user row having the given address.
// If the user does not exist yet, returns nil instead.
func (db DesmosDb) GetUserByAddress(address sdk.AccAddress) (*dbtypes.ProfileRow, error) {
	return db.executeQueryAndGetFirstUserRow(`SELECT * FROM profile WHERE address = $1`, address.String())
}
