package database

import "github.com/desmos-labs/djuno/v2/types"

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
	_, err := db.Sql.Exec(stmt, relationship.Creator, relationship.Counterparty, relationship.SubspaceID, relationship.Height)
	return err
}

// DeleteRelationship allows to delete the relationship between the given sender and receiver on the specified subspace
func (db Db) DeleteRelationship(relationship types.Relationship) error {
	stmt := `
DELETE FROM profile_relationship 
WHERE sender_address = $1 AND receiver_address = $2 AND subspace = $3 AND height <= $4`
	_, err := db.Sql.Exec(stmt,
		relationship.Creator, relationship.Counterparty, relationship.SubspaceID, relationship.Height)
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
	_, err := db.Sql.Exec(stmt, block.Blocker, block.Blocked, block.Reason, block.SubspaceID, block.Height)
	return err
}

// RemoveBlockage allow to remove a previously saved user blockage
func (db Db) RemoveBlockage(block types.Blockage) error {
	stmt := `
DELETE FROM user_block 
WHERE blocker_address = $1 AND blocked_user_address = $2 AND subspace = $3 AND height <= $4`
	_, err := db.Sql.Exec(stmt, block.Blocker, block.Blocked, block.SubspaceID, block.Height)
	return err
}
