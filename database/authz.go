package database

import (
	"fmt"
	"time"

	"github.com/desmos-labs/athena/types"
)

// SaveAuthzGrant saves the given grant inside the database
func (db *Db) SaveAuthzGrant(grant types.AuthzGrant) error {
	// Serialize the authorizations
	authzBz, err := db.cdc.MarshalInterfaceJSON(grant.Authorization)
	if err != nil {
		return fmt.Errorf("error while marshalling authorization to json: %s", err)
	}

	stmt := `
INSERT INTO authz_grant (granter_address, grantee_address, msg_type_url, "authorization", expiration, height)
VALUES ($1, $2, $3, $4, $5, $6) 
ON CONFLICT ON CONSTRAINT unique_msg_type_authorization DO UPDATE
    SET granter_address = excluded.granter_address,
        grantee_address = excluded.grantee_address,
        msg_type_url = excluded.msg_type_url,
        "authorization" = excluded."authorization",
        expiration = excluded.expiration,
        height = excluded.height
WHERE authz_grant.height <= excluded.height`
	_, err = db.SQL.Exec(stmt, grant.Granter, grant.Grantee, grant.Authorization.MsgTypeURL(), string(authzBz), grant.Expiration, grant.Height)
	return err
}

// DeleteAuthzGrant deletes the authz grant related to the given data
func (db *Db) DeleteAuthzGrant(granter string, grantee string, msgTypeURL string, height int64) error {
	stmt := `DELETE FROM authz_grant WHERE granter_address = $1 AND grantee_address = $2 AND msg_type_url = $3 AND height <= $4`
	_, err := db.SQL.Exec(stmt, granter, grantee, msgTypeURL, height)
	return err
}

// DeleteExpiredGrants deletes all the authz grants that are expired before or on the provided date
func (db *Db) DeleteExpiredGrants(time time.Time) error {
	stmt := `DELETE FROM authz_grant WHERE expiration <= $1`
	_, err := db.SQL.Exec(stmt, time)
	return err
}
