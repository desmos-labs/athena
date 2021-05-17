package database

import (
	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"

	"github.com/desmos-labs/djuno/x/profiles/types"

	dbtypes "github.com/desmos-labs/djuno/database/types"
)

// SaveUserIfNotExisting creates a new user having the given address if it does not exist yet.
// Upon creating the user, returns that.
// If any error is raised during the process, returns that.
func (db DesmosDb) SaveUserIfNotExisting(address string, height int64) error {
	stmt := `
INSERT INTO profile (address, height) 
VALUES ($1, $2)
ON CONFLICT (address) DO UPDATE 
    SET height = excluded.height
WHERE profile.height <= excluded.height`
	_, err := db.Sqlx.Exec(stmt, address, height)
	return err
}

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

	return dbtypes.ConvertProfileRow(rows[0])
}

// ---------------------------------------------------------------------------------------------------

// SaveProfile saves the given profile into the database, replacing any existing info.
// Returns the inserted row or an error if something goes wrong.
func (db DesmosDb) SaveProfile(profile types.Profile) error {
	stmt := `
INSERT INTO profile (address, nickname, dtag, bio, profile_pic, cover_pic, creation_time, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
ON CONFLICT (address) DO UPDATE 
	SET address = excluded.address, 
		nickname = excluded.nickname, 
		dtag = excluded.dtag,
		bio = excluded.bio,
		profile_pic = excluded.profile_pic,
		cover_pic = excluded.cover_pic,
		creation_time = excluded.creation_time,
		height = excluded.height
WHERE profile.height <= excluded.height`

	_, err := db.Sql.Exec(
		stmt,
		profile.GetAddress().String(), profile.Nickname, profile.DTag, profile.Bio,
		profile.Pictures.Profile, profile.Pictures.Cover, profile.CreationDate,
		profile.Height,
	)
	return err
}

// DeleteProfile allows to delete the profile of the user having the given address
func (db DesmosDb) DeleteProfile(address string, height int64) error {
	stmt := `
UPDATE profile 
SET nickname = '', 
    dtag = '', 
    bio = '', 
    profile_pic = '', 
    cover_pic = '' 
WHERE address = $1 AND height <= $2`
	_, err := db.Sql.Exec(stmt, address, height)
	return err
}

// ---------------------------------------------------------------------------------------------------

// SaveDTagTransferRequest saves a new transfer request from sender to receiver
func (db DesmosDb) SaveDTagTransferRequest(request types.DTagTransferRequest) error {
	err := db.SaveUserIfNotExisting(request.Sender, request.Height)
	if err != nil {
		return err
	}

	err = db.SaveUserIfNotExisting(request.Receiver, request.Height)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO dtag_transfer_requests (sender_address, receiver_address, height) 
VALUES ($1, $2, $3) 
ON CONFLICT ON CONSTRAINT unique_request DO UPDATE 
    SET sender_address = excluded.sender_address,
    	receiver_address = excluded.receiver_address
WHERE dtag_transfer_requests.height <= excluded.height`

	_, err = db.Sql.Exec(stmt, request.Sender, request.Receiver, request.Height)
	return err
}

// TransferDTag transfers the DTag from the sender to the receiver, and sets the sender DTag to the new one provided
func (db DesmosDb) TransferDTag(acceptance types.DTagTransferRequestAcceptance) error {
	// Get the old DTag
	var oldDTag string
	stmt := `SELECT dtag FROM profile WHERE address = $1 AND height <= $2`
	err := db.Sql.QueryRow(stmt, acceptance.Receiver, acceptance.Height).Scan(&oldDTag)
	if err != nil {
		return err
	}

	// Save the new DTags
	_, err = db.Sql.Exec(`UPDATE profile SET dtag = $1 WHERE address = $2`, acceptance.NewDTag, acceptance.Receiver)
	if err != nil {
		return err
	}

	_, err = db.Sql.Exec(`UPDATE profile SET dtag = $1 WHERE address = $2`, oldDTag, acceptance.Sender)
	if err != nil {
		return err
	}

	// Delete the transfer request
	return db.DeleteDTagTransferRequest(acceptance.DTagTransferRequest)
}

// DeleteDTagTransferRequest deletes the DTag requests from sender to receiver
func (db DesmosDb) DeleteDTagTransferRequest(request types.DTagTransferRequest) error {
	stmt := `
DELETE FROM dtag_transfer_requests 
WHERE sender_address = $1 AND receiver_address = $2 AND height <= $3`
	_, err := db.Sql.Exec(stmt, request.Sender, request.Receiver, request.Height)
	return err
}

// ---------------------------------------------------------------------------------------------------

// SaveRelationship allows to save a relationship between the sender and receiver on the given subspace
func (db DesmosDb) SaveRelationship(relationship types.Relationship) error {
	err := db.SaveUserIfNotExisting(relationship.Creator, relationship.Height)
	if err != nil {
		return err
	}

	err = db.SaveUserIfNotExisting(relationship.Recipient, relationship.Height)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO relationship (sender_address, receiver_address, subspace, height) 
VALUES ($1, $2, $3, $4) 
ON CONFLICT ON CONSTRAINT unique_relationship DO UPDATE 
    SET sender_address = excluded.sender_address,
		receiver_address = excluded.receiver_address,
		subspace = excluded.subspace
WHERE relationship.height <= excluded.height`
	_, err = db.Sql.Exec(stmt, relationship.Creator, relationship.Recipient, relationship.Subspace, relationship.Height)
	return err
}

// DeleteRelationship allows to delete the relationship between the given sender and receiver on the specified subspace
func (db DesmosDb) DeleteRelationship(relationship types.Relationship) error {
	stmt := `
DELETE FROM relationship 
WHERE sender_address = $1 AND receiver_address = $2 AND subspace = $3 AND height <= $4`
	_, err := db.Sql.Exec(stmt,
		relationship.Creator, relationship.Recipient, relationship.Subspace, relationship.Height)
	return err
}

// SaveBlockage allows to save a user blockage
func (db DesmosDb) SaveBlockage(block types.Blockage) error {
	err := db.SaveUserIfNotExisting(block.Blocker, block.Height)
	if err != nil {
		return err
	}

	err = db.SaveUserIfNotExisting(block.Blocker, block.Height)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO user_block(blocker_address, blocked_user_address, reason, subspace, height) 
VALUES ($1, $2, $3, $4, $5) 
ON CONFLICT ON CONSTRAINT unique_blockage DO UPDATE 
    SET blocker_address = excluded.blocker_address,
    	blocked_user_address = excluded.blocked_user_address,
    	reason = excluded.reason, 
    	subspace = excluded.subspace
WHERE user_block.height <= excluded.height`
	_, err = db.Sql.Exec(stmt, block.Blocker, block.Blocked, block.Reason, block.Subspace, block.Height)
	return err
}

// RemoveBlockage allow to remove a previously saved user blockage
func (db DesmosDb) RemoveBlockage(block types.Blockage) error {
	stmt := `
DELETE FROM user_block 
WHERE blocker_address = $1 AND blocked_user_address = $2 AND subspace = $3 AND height <= $4`
	_, err := db.Sql.Exec(stmt, block.Blocker, block.Blocked, block.Subspace, block.Height)
	return err
}
