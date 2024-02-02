package database

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	profilestypes "github.com/desmos-labs/desmos/v6/x/profiles/types"

	"github.com/desmos-labs/athena/types"

	dbtypes "github.com/desmos-labs/athena/database/types"
)

// SaveProfilesParams allows to store the given profiles params
func (db *Db) SaveProfilesParams(params types.ProfilesParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling profiles params: %s", err)
	}

	stmt := `
INSERT INTO profiles_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE profiles_params.height <= excluded.height`

	_, err = db.SQL.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing profiles params: %s", err)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------

// SaveUserIfNotExisting creates a new user having the given address if it does not exist yet.
// Upon creating the user, returns that.
// If any error is raised during the process, returns that.
func (db *Db) SaveUserIfNotExisting(address string, height int64) error {
	stmt := `INSERT INTO profile (address, height) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := db.SQL.Exec(stmt, address, height)
	return err
}

// GetUserByAddress returns the user row having the given address.
// If the user does not exist yet, returns nil instead.
func (db *Db) GetUserByAddress(address string) (*profilestypes.Profile, error) {
	var rows []dbtypes.ProfileRow
	err := db.SQL.Select(&rows, `SELECT * FROM profile WHERE address = $1`, address)
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
func (db *Db) SaveProfile(profile *types.Profile) error {
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

	_, err := db.SQL.Exec(
		stmt,
		profile.GetAddress().String(), profile.Nickname, profile.DTag, profile.Bio,
		profile.Pictures.Profile, profile.Pictures.Cover, profile.CreationDate,
		profile.Height,
	)
	return err
}

// DeleteProfile allows to delete the profile of the user having the given address
func (db *Db) DeleteProfile(address string, height int64) error {
	stmt := `DELETE FROM profile WHERE address = $1 AND height <= $2`
	_, err := db.SQL.Exec(stmt, address, height)
	return err
}

// GetProfilesAddresses returns all the addresses of the various profiles accounts
func (db *Db) GetProfilesAddresses() ([]string, error) {
	var rows []string
	err := db.SQL.Select(&rows, `SELECT address FROM profile`)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// ---------------------------------------------------------------------------------------------------

// SaveDTagTransferRequest saves a new transfer request from sender to receiver
func (db *Db) SaveDTagTransferRequest(request types.DTagTransferRequest) error {
	stmt := `
INSERT INTO dtag_transfer_requests (sender_address, receiver_address, height) 
VALUES ($1, $2, $3) 
ON CONFLICT ON CONSTRAINT unique_request DO UPDATE 
    SET sender_address = excluded.sender_address,
    	receiver_address = excluded.receiver_address
WHERE dtag_transfer_requests.height <= excluded.height`

	_, err := db.SQL.Exec(stmt, request.Sender, request.Receiver, request.Height)
	return err
}

// DeleteDTagTransferRequest deletes the DTag requests from sender to receiver
func (db *Db) DeleteDTagTransferRequest(request types.DTagTransferRequest) error {
	stmt := `
DELETE FROM dtag_transfer_requests 
WHERE sender_address = $1 AND receiver_address = $2 AND height <= $3`
	_, err := db.SQL.Exec(stmt, request.Sender, request.Receiver, request.Height)
	return err
}

// ---------------------------------------------------------------------------------------------------

// SaveChainLink allows to store inside the db the provided chain link
func (db *Db) SaveChainLink(link types.ChainLink) error {
	// Use a single transaction for the whole process
	tx, err := db.SQL.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Unpack the address data
	var address profilestypes.AddressData
	err = db.cdc.UnpackAny(link.Address, &address)
	if err != nil {
		return fmt.Errorf("error while reading link address as AddressData: %s", err)
	}

	// Insert the chain config
	chainConfigID, err := db.saveChainLinkChainConfig(tx, link.ChainConfig)
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
        creation_time = excluded.creation_time,
        height = excluded.height
WHERE chain_link.height <= excluded.height
RETURNING id`

	var chainLinkID int64
	err = tx.QueryRow(stmt, link.User, address.GetValue(), chainConfigID, link.CreationTime, link.Height).Scan(&chainLinkID)
	if err != nil {
		return err
	}

	// Insert the proof
	err = db.saveChainLinkProof(tx, chainLinkID, link.Proof, link.Height)
	if err != nil {
		return err
	}

	// Update the chain links count of the user
	stmt = `
INSERT INTO profile_counters (profile_address, chain_links_count)
VALUES ($1, 1)
ON CONFLICT (profile_address)
DO UPDATE SET chain_links_count = profile_counters.chain_links_count + 1;
`
	_, err = tx.Exec(stmt, link.User)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// getChainLinkID returns the id of the chain link for the given user, with the given chain config and external address
func (db *Db) getChainLinkID(userAddress string, chainConfigID int64, externalAddress string) (int64, error) {
	stmt := `SELECT id from chain_link WHERE user_address = $1 AND chain_config_id = $2 AND external_address = $3`

	var id int64
	err := db.SQL.QueryRow(stmt, userAddress, chainConfigID, externalAddress).Scan(&id)
	return id, err
}

// saveChainLinkProof stores the given proof as associated with the chain link having the given id
func (db *Db) saveChainLinkProof(tx *sql.Tx, chainLinkID int64, proof profilestypes.Proof, height int64) error {
	publicKeyBz, err := db.cdc.MarshalJSON(proof.PubKey)
	if err != nil {
		return fmt.Errorf("error serializing chain link proof public key: %s", err)
	}

	// Convert the plain text to hex if it's not already (might happen on old links)
	var plainText = proof.PlainText
	_, err = hex.DecodeString(plainText)
	if err != nil {
		plainText = hex.EncodeToString([]byte(plainText))
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

	signatureBz, err := db.cdc.MarshalJSON(proof.Signature)
	if err != nil {
		return fmt.Errorf("error serializing chain link signature: %s", err)
	}

	_, err = tx.Exec(stmt, chainLinkID, string(publicKeyBz), plainText, string(signatureBz), height)
	return err
}

// saveChainLinkChainConfig stores the given chain config and returns the row id
func (db *Db) saveChainLinkChainConfig(tx *sql.Tx, config profilestypes.ChainConfig) (int64, error) {
	stmt := `
INSERT INTO chain_link_chain_config (name) 
VALUES ($1)
ON CONFLICT ON CONSTRAINT unique_chain_config DO UPDATE 
    SET name = excluded.name
RETURNING id`

	var id int64
	err := tx.QueryRow(stmt, config.Name).Scan(&id)
	return id, err
}

// getChainLinkConfigID returns the chain link config id with the given name
func (db *Db) getChainLinkConfigID(name string) (int64, error) {
	stmt := `SELECT id FROM chain_link_chain_config WHERE name = $1`

	var id int64
	err := db.SQL.QueryRow(stmt, name).Scan(&id)
	return id, err
}

// DeleteChainLink removes from the database the chain link made for the given user and having the provided
// external address and linked to the chain with the given name
func (db *Db) DeleteChainLink(user string, externalAddress string, chainName string, height int64) error {
	tx, err := db.SQL.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete the chain link
	stmt := `
DELETE FROM chain_link 
WHERE user_address = $1 
  AND external_address = $2
  AND chain_config_id = (SELECT id FROM chain_link_chain_config WHERE name = $3)
  AND height <= $4`
	_, err = tx.Exec(stmt, user, externalAddress, chainName, height)
	if err != nil {
		return err
	}

	// Update the chain links count of the user
	stmt = `
INSERT INTO profile_counters (profile_address, chain_links_count)
VALUES ($1, 0)
ON CONFLICT (profile_address)
DO UPDATE SET chain_links_count = profile_counters.chain_links_count - 1;
`
	_, err = tx.Exec(stmt, user)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteAllChainLinks deletes all the chain links having a height lower than the given one
func (db *Db) DeleteAllChainLinks(height int64) error {
	tx, err := db.SQL.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete the chain links
	stmt := `DELETE FROM chain_link WHERE height <= $1`
	_, err = tx.Exec(stmt, height)
	if err != nil {
		return err
	}

	// Reset the chain links count of all the users
	stmt = `UPDATE profile_counters SET chain_links_count = 0 WHERE TRUE`
	_, err = tx.Exec(stmt)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// SaveDefaultChainLink saves the given chain link as a default chain link
func (db *Db) SaveDefaultChainLink(chainLink types.ChainLink) error {
	stmt := `
INSERT INTO default_chain_link (user_address, chain_link_id, chain_config_id, height) 
VALUES ($1, $2, $3, $4)
ON CONFLICT ON CONSTRAINT unique_default_chain_link DO UPDATE 
    SET chain_link_id = excluded.chain_link_id, 
        chain_config_id = excluded.chain_config_id, 
        height = excluded.height
WHERE default_chain_link.height <= excluded.height`

	chainLinkConfigID, err := db.getChainLinkConfigID(chainLink.ChainConfig.Name)
	if err != nil {
		return err
	}

	var address profilestypes.AddressData
	err = db.cdc.UnpackAny(chainLink.Address, &address)
	if err != nil {
		return fmt.Errorf("error while reading link address as AddressData: %s", err)
	}

	chainLinkID, err := db.getChainLinkID(chainLink.User, chainLinkConfigID, address.GetValue())
	if err != nil {
		return err
	}

	_, err = db.SQL.Exec(stmt, chainLink.User, chainLinkID, chainLinkConfigID, chainLink.Height)
	return err
}

// DeleteAllDefaultChainLinks removes all default chain links having a height lower than the one specified
func (db *Db) DeleteAllDefaultChainLinks(height int64) error {
	stmt := `DELETE FROM default_chain_link WHERE height <= $1`
	_, err := db.SQL.Exec(stmt, height)
	return err
}

// ---------------------------------------------------------------------------------------------------

// getApplicationLinkRowID returns the row id of the application link for
// the provided application, username and Desmos address
func (db *Db) getApplicationLinkRowID(address string, application string, username string) (sql.NullInt64, error) {
	stmt := `SELECT id FROM application_link WHERE user_address = $1 AND application ILIKE $2 AND username ILIKE $3`

	var rowID int64
	err := db.SQL.QueryRow(stmt, address, application, username).Scan(&rowID)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.NullInt64{Int64: 0, Valid: false}, nil
	}

	return sql.NullInt64{Int64: rowID, Valid: true}, err
}

// SaveApplicationLink stores the given application link inside the database
func (db *Db) SaveApplicationLink(link types.ApplicationLink) error {
	// Use a single transaction for the whole process
	tx, err := db.SQL.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Save the link
	stmt := `
INSERT INTO application_link (user_address, application, username, state, result, creation_time, expiration_time, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT ON CONSTRAINT unique_application_link DO UPDATE 
    SET user_address = excluded.user_address, 
    	application = excluded.application,
    	username = excluded.username,
    	state = excluded.state,
    	result = excluded.result, 
    	creation_time = excluded.creation_time,
    	expiration_time = excluded.expiration_time,
    	height = excluded.height
WHERE application_link.height <= excluded.height
RETURNING id`

	var result sql.NullString
	if link.Result != nil {
		resultBz, err := db.cdc.MarshalJSON(link.Result)
		if err != nil {
			return fmt.Errorf("error while serializing result: %s", err)
		}
		result = sql.NullString{Valid: true, String: string(resultBz)}
	}

	var linkID int64
	err = tx.QueryRow(stmt,
		link.User,
		link.Data.Application,
		link.Data.Username,
		link.State.String(),
		result,
		link.CreationTime,
		link.ExpirationTime,
		link.Height,
	).Scan(&linkID)
	if err != nil {
		return err
	}

	// Save the oracle request
	err = db.saveOracleRequest(tx, linkID, link.OracleRequest, link.Height)
	if err != nil {
		return err
	}

	// Update the application links count of the user
	stmt = `
INSERT INTO profile_counters (profile_address, application_links_count)
VALUES ($1, 1)
ON CONFLICT (profile_address)
DO UPDATE SET application_links_count = profile_counters.application_links_count + 1;
`
	_, err = tx.Exec(stmt, link.User)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// saveOracleRequest stores the given oracle request associating it with the link having the provided id
func (db *Db) saveOracleRequest(tx *sql.Tx, linkID int64, request profilestypes.OracleRequest, height int64) error {
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

	callDataBz, err := db.cdc.MarshalJSON(&request.CallData)
	if err != nil {
		return fmt.Errorf("error while serializing oracle request call data: %s", err)
	}

	_, err = tx.Exec(stmt,
		linkID,
		fmt.Sprintf("%d", request.ID),
		fmt.Sprintf("%d", request.OracleScriptID),
		string(callDataBz),
		request.ClientID,
		height,
	)
	return err
}

type applicationLinkInfo struct {
	User        string `db:"user_address"`
	Application string `db:"application"`
	Username    string `db:"username"`
}

// GetApplicationLinkInfos returns the information of all the stored application links
func (db *Db) GetApplicationLinkInfos() ([]types.ApplicationLinkInfo, error) {
	stmt := `SELECT user_address, application, username FROM application_link`

	var rows []applicationLinkInfo
	err := db.SQL.Select(&rows, stmt)
	if err != nil {
		return nil, err
	}

	applicationLinkInfos := make([]types.ApplicationLinkInfo, len(rows))
	for i, row := range rows {
		applicationLinkInfos[i] = types.NewApplicationInfo(row.User, row.Application, row.Username)
	}
	return applicationLinkInfos, nil
}

// DeleteApplicationLink allows to delete the application link associated to the given user,
// having the given application and username values
func (db *Db) DeleteApplicationLink(user, application, username string, height int64) error {
	tx, err := db.SQL.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete the link
	stmt := `
DELETE FROM application_link 
WHERE user_address = $1 
  AND application = $2 
  AND username = $3 
  AND height <= $4`
	_, err = tx.Exec(stmt, user, application, username, height)
	if err != nil {
		return err
	}

	// Update the application links count of the user
	stmt = `
	INSERT INTO profile_counters (profile_address, application_links_count)
VALUES ($1, 0)
ON CONFLICT (profile_address)
DO UPDATE SET application_links_count = profile_counters.application_links_count - 1;
`
	_, err = tx.Exec(stmt, user)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteAllApplicationLinks deletes all the application links that have a height equal or lower to the one given
func (db *Db) DeleteAllApplicationLinks(height int64) error {
	tx, err := db.SQL.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete the application links
	stmt := `DELETE FROM application_link WHERE height <= $1`
	_, err = tx.Exec(stmt, height)
	if err != nil {
		return err
	}

	// Reset the application links count of all the users
	stmt = `UPDATE profile_counters SET application_links_count = 0 WHERE TRUE`
	_, err = tx.Exec(stmt)
	if err != nil {
		return err
	}

	return tx.Commit()
}
