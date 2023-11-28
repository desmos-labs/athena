package database

import (
	"fmt"

	"github.com/lib/pq"

	dbtypes "github.com/desmos-labs/athena/database/types"
	"github.com/desmos-labs/athena/types"
)

// SaveTip saves the given tip inside the database
func (db *Db) SaveTip(tip types.Tip) error {
	switch tipTarget := tip.Target.(type) {
	case types.UserTarget:
		return db.saveUserTip(tip, tipTarget)
	case types.PostTarget:
		return db.savePostTip(tip, tipTarget)

	default:
		return fmt.Errorf("invalid tip target: %T", tip.Target)
	}
}

// saveUserTip saves the given tip for the provided target inside the database
func (db *Db) saveUserTip(tip types.Tip, target types.UserTarget) error {
	stmt := `
INSERT INTO tip_user (sender_address, receiver_address, subspace_id, amount, height) 
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT ON CONSTRAINT unique_sender_user_tip DO UPDATE 
    SET sender_address = excluded.sender_address,
        receiver_address = excluded.receiver_address,
        subspace_id = excluded.subspace_id, 
        amount = excluded.amount,
        height = excluded.height
WHERE tip_user.height <= excluded.height`
	_, err := db.SQL.Exec(stmt, tip.Sender, target.Address, tip.SubspaceID, pq.Array(dbtypes.NewDbCoins(tip.Amount)), tip.Height)
	return err
}

// savePostTip saves the given tip for the provided target inside the database
func (db *Db) savePostTip(tip types.Tip, target types.PostTarget) error {
	postRowID, err := db.getPostRowID(tip.SubspaceID, target.PostID)
	if err != nil {
		return fmt.Errorf("no row found for post %d within subspace %d", target.PostID, tip.SubspaceID)
	}

	stmt := `
INSERT INTO tip_post (sender_address, subspace_id, post_row_id, amount, height) 
VALUES ($1, $2, $3, $4, $5) 
ON CONFLICT ON CONSTRAINT unique_sender_post_tip DO UPDATE 
    SET sender_address = excluded.sender_address,
        subspace_id = excluded.subspace_id,
        post_row_id = excluded.post_row_id,
        amount = excluded.amount,
        height = excluded.height
WHERE tip_post.height <= excluded.height`
	_, err = db.SQL.Exec(stmt, tip.Sender, tip.SubspaceID, postRowID, pq.Array(dbtypes.NewDbCoins(tip.Amount)), tip.Height)
	return err
}
