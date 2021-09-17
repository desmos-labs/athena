package database

import (
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"

	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"

	"github.com/desmos-labs/djuno/types"

	dbtypes "github.com/desmos-labs/djuno/database/types"
)

// SaveUserIfNotExisting creates a new user having the given address if it does not exist yet.
// Upon creating the user, returns that.
// If any error is raised during the process, returns that.
func (db Db) SaveUserIfNotExisting(address string, height int64) error {
	stmt := `INSERT INTO profile (address, height) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := db.Sqlx.Exec(stmt, address, height)
	return err
}

// GetUserByAddress returns the user row having the given address.
// If the user does not exist yet, returns nil instead.
func (db Db) GetUserByAddress(address string) (*profilestypes.Profile, error) {
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
func (db Db) SaveProfile(profile types.Profile) error {
	log.Info().Str("dtag", profile.DTag).Msg("saving profile")

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
func (db Db) DeleteProfile(address string, height int64) error {
	log.Info().Str("address", address).Msg("deleting profile")

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
func (db Db) SaveDTagTransferRequest(request types.DTagTransferRequest) error {
	stmt := `
INSERT INTO dtag_transfer_requests (sender_address, receiver_address, height) 
VALUES ($1, $2, $3) 
ON CONFLICT ON CONSTRAINT unique_request DO UPDATE 
    SET sender_address = excluded.sender_address,
    	receiver_address = excluded.receiver_address
WHERE dtag_transfer_requests.height <= excluded.height`

	_, err := db.Sql.Exec(stmt, request.Sender, request.Receiver, request.Height)
	return err
}

// DeleteDTagTransferRequest deletes the DTag requests from sender to receiver
func (db Db) DeleteDTagTransferRequest(request types.DTagTransferRequest) error {
	stmt := `
DELETE FROM dtag_transfer_requests 
WHERE sender_address = $1 AND receiver_address = $2 AND height <= $3`
	_, err := db.Sql.Exec(stmt, request.Sender, request.Receiver, request.Height)
	return err
}

// ---------------------------------------------------------------------------------------------------

// SaveRelationship allows to save a relationship between the sender and receiver on the given subspace
func (db Db) SaveRelationship(relationship types.Relationship) error {
	stmt := `
INSERT INTO profile_relationship (sender_address, receiver_address, subspace, height) 
VALUES ($1, $2, $3, $4) 
ON CONFLICT ON CONSTRAINT unique_relationship DO UPDATE 
    SET sender_address = excluded.sender_address,
		receiver_address = excluded.receiver_address,
		subspace = excluded.subspace
WHERE profile_relationship.height <= excluded.height`
	_, err := db.Sql.Exec(stmt, relationship.Creator, relationship.Recipient, relationship.Subspace, relationship.Height)
	return err
}

// DeleteRelationship allows to delete the relationship between the given sender and receiver on the specified subspace
func (db Db) DeleteRelationship(relationship types.Relationship) error {
	stmt := `
DELETE FROM profile_relationship 
WHERE sender_address = $1 AND receiver_address = $2 AND subspace = $3 AND height <= $4`
	_, err := db.Sql.Exec(stmt,
		relationship.Creator, relationship.Recipient, relationship.Subspace, relationship.Height)
	return err
}

// ---------------------------------------------------------------------------------------------------

// SaveBlockage allows to save a user blockage
func (db Db) SaveBlockage(block types.Blockage) error {
	stmt := `
INSERT INTO user_block(blocker_address, blocked_user_address, reason, subspace, height) 
VALUES ($1, $2, $3, $4, $5) 
ON CONFLICT ON CONSTRAINT unique_blockage DO UPDATE 
    SET blocker_address = excluded.blocker_address,
    	blocked_user_address = excluded.blocked_user_address,
    	reason = excluded.reason, 
    	subspace = excluded.subspace
WHERE user_block.height <= excluded.height`
	_, err := db.Sql.Exec(stmt, block.Blocker, block.Blocked, block.Reason, block.Subspace, block.Height)
	return err
}

// RemoveBlockage allow to remove a previously saved user blockage
func (db Db) RemoveBlockage(block types.Blockage) error {
	stmt := `
DELETE FROM user_block 
WHERE blocker_address = $1 AND blocked_user_address = $2 AND subspace = $3 AND height <= $4`
	_, err := db.Sql.Exec(stmt, block.Blocker, block.Blocked, block.Subspace, block.Height)
	return err
}

// ---------------------------------------------------------------------------------------------------

// SaveChainLink allows to store inside the db the provided chain link
func (db Db) SaveChainLink(link types.ChainLink) error {
	// Insert the chain config
	chainConfigID, err := db.saveChainLinkChainConfig(link.ChainConfig)
	if err != nil {
		return err
	}

	// Insert the link details and get the id
	stmt := `
INSERT INTO chain_link (user_address, external_address, chain_config_id, creation_time, height)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT ON CONSTRAINT unique_chain_link DO UPDATE 
    SET user_address = excluded.user_address, 
        external_address = excluded.external_address,
        chain_config_id = excluded.chain_config_id,
        creation_time = excluded.creation_time,
        height = excluded.height
WHERE chain_link.height <= excluded.height
RETURNING id`

	var address profilestypes.AddressData
	err = db.EncodingConfig.Marshaler.UnpackAny(link.Address, &address)
	if err != nil {
		return fmt.Errorf("error while reading link address as AddressData: %s", err)
	}

	var chainLinkID int64
	err = db.Sql.
		QueryRow(stmt, link.User, address.GetValue(), chainConfigID, link.CreationTime, link.Height).
		Scan(&chainLinkID)
	if err != nil {
		return err
	}

	// Insert the proof
	return db.saveChainLinkProof(chainLinkID, link.Proof, link.Height)
}

// saveChainLinkProof stores the given proof as associated with the chain link having the given id
func (db Db) saveChainLinkProof(chainLinkID int64, proof profilestypes.Proof, height int64) error {
	publicKeyBz, err := db.EncodingConfig.Marshaler.MarshalJSON(proof.PubKey)
	if err != nil {
		return fmt.Errorf("error serializing chain link proof public key: %s", err)
	}

	stmt := `
INSERT INTO chain_link_proof(chain_link_id, public_key, plain_text, signature, height) 
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT ON CONSTRAINT unique_proof_for_link DO UPDATE
    SET chain_link_id = excluded.chain_link_id, 
        public_key = excluded.public_key, 
        plain_text = excluded.plain_text, 
        signature = excluded.signature, 
        height = excluded.height
WHERE chain_link_proof.height <= excluded.height`
	_, err = db.Sql.Exec(stmt, chainLinkID, string(publicKeyBz), proof.PlainText, proof.Signature, height)
	return err
}

// saveChainLinkChainConfig stores the given chain config and returns the row id
func (db Db) saveChainLinkChainConfig(config profilestypes.ChainConfig) (int64, error) {
	stmt := `
INSERT INTO chain_link_chain_config (name) 
VALUES ($1)
ON CONFLICT ON CONSTRAINT unique_chain_config DO UPDATE 
    SET name = excluded.name
RETURNING id`

	var id int64
	err := db.Sql.QueryRow(stmt, config.Name).Scan(&id)
	return id, err
}

// DeleteChainLink removes from the database the chain link made for the given user and having the provided
// external address and linked to the chain with the given name
func (db Db) DeleteChainLink(user string, externalAddress string, chainName string) error {
	stmt := `
DELETE FROM chain_link 
WHERE user_address = $1 
  AND external_address = $2
  AND chain_config_id = (SELECT id FROM chain_link_chain_config WHERE name = $3)`
	_, err := db.Sql.Exec(stmt, user, externalAddress, chainName)
	return err
}

// ---------------------------------------------------------------------------------------------------

// SaveApplicationLink stores the given application link inside the database
func (db Db) SaveApplicationLink(link types.ApplicationLink) error {
	// Save the link
	stmt := `
INSERT INTO application_link (user_address, application, username, state, result, creation_time, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT ON CONSTRAINT unique_application_link DO UPDATE 
    SET user_address = excluded.user_address, 
    	application = excluded.application,
    	username = excluded.username,
    	state = excluded.state,
    	result = excluded.result, 
    	creation_time = excluded.creation_time,
    	height = excluded.height
WHERE application_link.height <= excluded.height
RETURNING id`

	var result sql.NullString
	if link.Result != nil {
		resultBz, err := db.EncodingConfig.Marshaler.MarshalJSON(link.Result)
		if err != nil {
			return fmt.Errorf("error while serializing result: %s", err)
		}
		result = sql.NullString{Valid: true, String: string(resultBz)}
	}

	var linkID int64
	err := db.Sql.QueryRow(stmt,
		link.User, link.Data.Application, link.Data.Username, link.State.String(),
		result, link.CreationTime, link.Height,
	).Scan(&linkID)
	if err != nil {
		return err
	}

	// Save the oracle request
	return db.saveOracleRequest(linkID, link.OracleRequest, link.Height)
}

// saveOracleRequest stores the given oracle request associating it with the link having the provided id
func (db Db) saveOracleRequest(linkID int64, request profilestypes.OracleRequest, height int64) error {
	stmt := `
INSERT INTO application_link_oracle_request (application_link_id, request_id, script_id, call_data, client_id, height) 
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT ON CONSTRAINT unique_oracle_request DO UPDATE 
    SET application_link_id = excluded.application_link_id,
        request_id = excluded.request_id, 
        script_id = excluded.script_id, 
        call_data = excluded.call_data, 
        client_id = excluded.client_id,
        height = excluded.height
WHERE application_link_oracle_request.height <= excluded.height`

	callDataBz, err := db.EncodingConfig.Marshaler.MarshalJSON(&request.CallData)
	if err != nil {
		return fmt.Errorf("error while serializing oracle request call data: %s", err)
	}

	_, err = db.Sql.Exec(stmt, linkID, request.ID, request.OracleScriptID, string(callDataBz), request.ClientID, height)
	return err
}

// DeleteApplicationLink allows to delete the application link associated to the given user,
// having the given application and username values
func (db Db) DeleteApplicationLink(user, application, username string) error {
	stmt := `DELETE FROM application_link WHERE user_address = $1 AND application = $2 AND username = $3`
	_, err := db.Sql.Exec(stmt, user, application, username)
	return err
}
