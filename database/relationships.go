package database

import (
	"github.com/desmos-labs/athena/v2/types"
)

// SaveRelationship allows to save a relationship between the sender and receiver on the given subspace
func (db *Db) SaveRelationship(relationship types.Relationship) error {
	tx, err := db.SQL.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Store the relationship
	stmt := `
INSERT INTO user_relationship (creator_address, counterparty_address, subspace_id, height) 
VALUES ($1, $2, $3, $4) 
ON CONFLICT ON CONSTRAINT unique_relationship DO UPDATE 
    SET creator_address = excluded.creator_address,
		counterparty_address = excluded.counterparty_address,
		subspace_id = excluded.subspace_id,
		height = excluded.height
WHERE user_relationship.height <= excluded.height`
	_, err = tx.Exec(stmt, relationship.Creator, relationship.Counterparty, relationship.SubspaceID, relationship.Height)
	if err != nil {
		return err
	}

	// Update the relationships count of the creator
	stmt = `
INSERT INTO profile_counters (profile_address, relationships_count)
VALUES ($1, 1)
ON CONFLICT (profile_address)
DO UPDATE SET relationships_count = profile_counters.relationships_count + 1;
`
	_, err = tx.Exec(stmt, relationship.Creator)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteRelationship allows to delete the relationship between the given sender and receiver on the specified subspace
func (db *Db) DeleteRelationship(relationship types.Relationship) error {
	tx, err := db.SQL.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete the relationship
	stmt := `
DELETE FROM user_relationship 
WHERE creator_address = $1 AND counterparty_address = $2 AND subspace_id = $3 AND height <= $4`
	_, err = tx.Exec(stmt, relationship.Creator, relationship.Counterparty, relationship.SubspaceID, relationship.Height)
	if err != nil {
		return err
	}

	// Update the relationships count of the creator
	stmt = `
INSERT INTO profile_counters (profile_address, relationships_count)
VALUES ($1, 0)
ON CONFLICT (profile_address)
DO UPDATE SET relationships_count = profile_counters.relationships_count - 1;
`
	_, err = tx.Exec(stmt, relationship.Creator)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteAllRelationships allows to delete all the relationships associated with the given subspace from the database
func (db *Db) DeleteAllRelationships(height int64, subspaceID uint64) error {
	tx, err := db.SQL.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete all the relationships
	stmt := `DELETE FROM user_relationship WHERE subspace_id = $1 AND height <= $2`
	_, err = tx.Exec(stmt, subspaceID, height)
	if err != nil {
		return err
	}

	// Delete all the relationships counters
	stmt = `UPDATE profile_counters SET relationships_count = 0 WHERE profile_address = $1`
	_, err = tx.Exec(stmt, subspaceID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// ---------------------------------------------------------------------------------------------------

// SaveUserBlock allows to save a user blockage
func (db *Db) SaveUserBlock(block types.Blockage) error {
	tx, err := db.SQL.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Store the user block
	stmt := `
INSERT INTO user_block (blocker_address, blocked_address, reason, subspace_id, height) 
VALUES ($1, $2, $3, $4, $5) 
ON CONFLICT ON CONSTRAINT unique_blockage DO UPDATE 
    SET blocker_address = excluded.blocker_address,
    	blocked_address = excluded.blocked_address,
    	reason = excluded.reason, 
    	subspace_id = excluded.subspace_id
WHERE user_block.height <= excluded.height`
	_, err = tx.Exec(stmt, block.Blocker, block.Blocked, block.Reason, block.SubspaceID, block.Height)
	if err != nil {
		return err
	}

	// Update the blocks count of the blocker
	stmt = `
INSERT INTO profile_counters (profile_address, blocks_count)
VALUES ($1, 1)
ON CONFLICT (profile_address)
DO UPDATE SET blocks_count = profile_counters.blocks_count + 1;
`
	_, err = tx.Exec(stmt, block.Blocker)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteBlockage allow to remove a previously saved user blockage
func (db *Db) DeleteBlockage(block types.Blockage) error {
	tx, err := db.SQL.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete the blockage
	stmt := `
DELETE FROM user_block 
WHERE blocker_address = $1 AND blocked_address = $2 AND subspace_id = $3 AND height <= $4`
	_, err = tx.Exec(stmt, block.Blocker, block.Blocked, block.SubspaceID, block.Height)
	if err != nil {
		return err
	}

	// Update the blocks count of the blocker
	stmt = `
INSERT INTO profile_counters (profile_address, blocks_count)
VALUES ($1, 0)
ON CONFLICT (profile_address)
DO UPDATE SET blocks_count = profile_counters.blocks_count - 1;
`
	_, err = tx.Exec(stmt, block.Blocker)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteAllUserBlocks allows to delete all the user blocks associated with the given subspace from the database
func (db *Db) DeleteAllUserBlocks(height int64, subspaceID uint64) error {
	tx, err := db.SQL.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete all the user blocks
	stmt := `DELETE FROM user_block WHERE subspace_id = $1 AND height <= $2`
	_, err = tx.Exec(stmt, subspaceID, height)
	if err != nil {
		return err
	}

	// Set the blocks counter to 0
	stmt = `UPDATE profile_counters SET blocks_count = 0 WHERE profile_address = $1`
	_, err = tx.Exec(stmt, subspaceID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
