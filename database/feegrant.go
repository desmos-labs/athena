package database

import (
	"fmt"
	"time"

	"github.com/cosmos/gogoproto/proto"
	"github.com/lib/pq"

	dbtypes "github.com/desmos-labs/djuno/v2/database/types"
	"github.com/desmos-labs/djuno/v2/types"
)

// SaveFeeGrant stores the given grant inside the database
func (db *Db) SaveFeeGrant(grant types.FeeGrant) error {
	// Serialize the allowance
	msg, ok := grant.Allowance.(proto.Message)
	if !ok {
		return fmt.Errorf("cannot proto marshal %T", grant.Allowance)
	}
	allowanceBz, err := db.cdc.MarshalInterfaceJSON(msg)
	if err != nil {
		return fmt.Errorf("cannot encode %T", msg)
	}

	// Get the spend limit
	spendLimit, err := grant.GetSpendLimit()
	if err != nil {
		return err
	}

	// Get the expiration time
	expirationTime, err := grant.GetExpirationDate()
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO fee_grant (granter_address, grantee_address, spend_limit, expiration_date, allowance, height) 
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT ON CONSTRAINT unique_fee_grant DO UPDATE 
    SET granter_address = excluded.granter_address,
        grantee_address = excluded.grantee_address,
        spend_limit = excluded.spend_limit,
        expiration_date = excluded.expiration_date,
        allowance = excluded.allowance,
        height = excluded.height
WHERE fee_grant.height <= excluded.height`

	_, err = db.SQL.Exec(stmt,
		grant.Granter,
		grant.Grantee,
		pq.Array(dbtypes.NewDbCoins(spendLimit)),
		expirationTime,
		string(allowanceBz),
		grant.Height,
	)
	return err
}

// DeleteFeeGrant removes the fee grant for the given data from the database
func (db *Db) DeleteFeeGrant(granter string, grantee string, height int64) error {
	stmt := `DELETE FROM fee_grant WHERE granter_address = $1 AND grantee_address = $2 AND height <= $3`
	_, err := db.SQL.Exec(stmt, granter, grantee, height)
	return err
}

// DeleteExpiredFeeGrants removes the fee grants that expire before or on the given time
func (db *Db) DeleteExpiredFeeGrants(time time.Time) error {
	stmt := `DELETE FROM fee_grant WHERE expiration_date <= $1`
	_, err := db.SQL.Exec(stmt, time)
	return err
}
