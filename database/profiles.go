package database

import (
	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"
	dbtypes "github.com/desmos-labs/djuno/database/types"
)

// SaveUserIfNotExisting creates a new user having the given address if it does not exist yet.
// Upon creating the user, returns that.
// If any error is raised during the process, returns that.
func (db DesmosDb) SaveUserIfNotExisting(address string) error {
	stmt := `INSERT INTO profile (address) VALUES ($1) ON CONFLICT DO NOTHING`
	_, err := db.Sqlx.Exec(stmt, address)
	return err
}

// SaveProfile saves the given profile into the database, replacing any existing info.
// Returns the inserted row or an error if something goes wrong.
func (db DesmosDb) SaveProfile(profile profilestypes.Profile) error {
	stmt := `
INSERT INTO profile (address, moniker, dtag, bio, profile_pic, cover_pic, creation_time) 
VALUES ($1, $2, $3, $4, $5, $6, $7) 
ON CONFLICT (address) DO UPDATE 
	SET address = excluded.address, 
		moniker = excluded.moniker, 
		dtag = excluded.dtag,
		bio = excluded.bio,
		profile_pic = excluded.profile_pic,
		cover_pic = excluded.cover_pic,
		creation_time = excluded.creation_time`

	_, err := db.Sql.Exec(
		stmt,
		profile.Creator, profile.Moniker, profile.Dtag, profile.Bio,
		profile.Pictures.Profile, profile.Pictures.Cover, profile.CreationDate,
	)
	return err
}

// ---------------------------------------------------------------------------------------------------

// GetUserByAddress returns the user row having the given address.
// If the user does not exist yet, returns nil instead.
func (db DesmosDb) GetUserByAddress(address string) (*profilestypes.Profile, error) {
	var rows []dbtypes.ProfileRow
	err := db.Sqlx.Select(&rows, `SELECT * FROM profile WHERE address = $1`, address)
	if err != nil {
		return nil, err
	}

	// No users found, return nil
	if len(rows) == 0 {
		return nil, nil
	}

	profile := dbtypes.ConvertProfileRow(rows[0])
	return &profile, nil
}

// ---------------------------------------------------------------------------------------------------

// DeleteProfile allows to delete the profile of the user having the given address
func (db DesmosDb) DeleteProfile(address string) error {
	stmt := `UPDATE profile SET moniker = '', dtag = '', bio = '', profile_pic = '', cover_pic = '' WHERE address = $1`
	_, err := db.Sql.Exec(stmt, address)
	return err
}

// ---------------------------------------------------------------------------------------------------

// SaveDTagTransferRequest saves a new transfer request from sender to receiver
func (db DesmosDb) SaveDTagTransferRequest(sender, receiver string) error {
	err := db.SaveUserIfNotExisting(sender)
	if err != nil {
		return err
	}

	err = db.SaveUserIfNotExisting(receiver)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO dtag_transfer_requests (sender_address, receiver_address) VALUES ($2, $3) ON CONFLICT DO NOTHING`

	_, err = db.Sql.Exec(stmt, sender, receiver)
	return err
}

func (db DesmosDb) TransferDTag(newDTag, sender, receiver string) error {
	// Get the old DTag
	var oldDTag string
	err := db.Sql.QueryRow(`SELECT dtag FROM profile WHERE address = $1`, receiver).Scan(&oldDTag)
	if err != nil {
		return err
	}

	// Save the new DTags
	_, err = db.Sql.Exec(`UPDATE profile SET dtag = $1 WHERE address = $2`, newDTag, receiver)
	if err != nil {
		return err
	}

	_, err = db.Sql.Exec(`UPDATE profile SET dtag = $1 WHERE address = $2`, oldDTag, sender)
	if err != nil {
		return err
	}

	// Delete the transfer request
	return db.DeleteDTagTransferRequest(sender, receiver)
}

func (db DesmosDb) DeleteDTagTransferRequest(sender, receiver string) error {
	_, err := db.Sql.Exec(`DELETE FROM dtag_transfer_requests WHERE sender_address = $1 AND receiver_address = $2`)
	return err
}
