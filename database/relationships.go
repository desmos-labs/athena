package database

import (
	"github.com/desmos-labs/djuno/v2/types"
)

// SaveRelationship allows to save a relationship between the sender and receiver on the given subspace
func (db *Db) SaveRelationship(relationship types.Relationship) error {
	stmt := `
INSERT INTO user_relationship (creator_address, counterparty_address, subspace_id, height) 
VALUES ($1, $2, $3, $4) 
ON CONFLICT ON CONSTRAINT unique_relationship DO UPDATE 
    SET creator_address = excluded.creator_address,
		counterparty_address = excluded.counterparty_address,
		subspace_id = excluded.subspace_id
WHERE user_relationship.height <= excluded.height`
	_, err := db.Sql.Exec(stmt, relationship.Creator, relationship.Counterparty, relationship.SubspaceID, relationship.Height)
	return err
}

// DeleteRelationship allows to delete the relationship between the given sender and receiver on the specified subspace
func (db *Db) DeleteRelationship(relationship types.Relationship) error {
	stmt := `
DELETE FROM user_relationship 
WHERE creator_address = $1 AND counterparty_address = $2 AND subspace_id = $3 AND height <= $4`
	_, err := db.Sql.Exec(stmt, relationship.Creator, relationship.Counterparty, relationship.SubspaceID, relationship.Height)
	return err
}

// DeleteAllRelationships allows to delete all the relationships associated with the given subspace from the database
func (db *Db) DeleteAllRelationships(height int64, subspaceID uint64) error {
	stmt := `DELETE FROM user_relationship WHERE subspace_id = $1 AND height <= $2`
	_, err := db.Sql.Exec(stmt, subspaceID, height)
	return err
}

// ---------------------------------------------------------------------------------------------------

// SaveUserBlock allows to save a user blockage
func (db *Db) SaveUserBlock(block types.Blockage) error {
	stmt := `
INSERT INTO user_block (blocker_address, blocked_address, reason, subspace_id, height) 
VALUES ($1, $2, $3, $4, $5) 
ON CONFLICT ON CONSTRAINT unique_blockage DO UPDATE 
    SET blocker_address = excluded.blocker_address,
    	blocked_address = excluded.blocked_address,
    	reason = excluded.reason, 
    	subspace_id = excluded.subspace_id
WHERE user_block.height <= excluded.height`
	_, err := db.Sql.Exec(stmt, block.Blocker, block.Blocked, block.Reason, block.SubspaceID, block.Height)
	return err
}

// DeleteBlockage allow to remove a previously saved user blockage
func (db *Db) DeleteBlockage(block types.Blockage) error {
	stmt := `
DELETE FROM user_block 
WHERE blocker_address = $1 AND blocked_address = $2 AND subspace_id = $3 AND height <= $4`
	_, err := db.Sql.Exec(stmt, block.Blocker, block.Blocked, block.SubspaceID, block.Height)
	return err
}

// DeleteAllUserBlocks allows to delete all the user blocks associated with the given subspace from the database
func (db *Db) DeleteAllUserBlocks(height int64, subspaceID uint64) error {
	stmt := `DELETE FROM user_block WHERE subspace_id = $1 AND height <= $2`
	_, err := db.Sql.Exec(stmt, subspaceID, height)
	return err
}
