package database

import (
	poststypes "github.com/desmos-labs/desmos/v2/x/staging/posts/types"

	"github.com/desmos-labs/djuno/types"

	dbtypes "github.com/desmos-labs/djuno/database/types"
)

// convertPostRow takes the given row and returns a RegisteredReaction
func convertReactionRow(row dbtypes.RegisteredReactionRow) poststypes.RegisteredReaction {
	return poststypes.NewRegisteredReaction(
		row.OwnerAddress,
		row.ShortCode,
		row.Value,
		row.Subspace,
	)
}

// SavePostReaction allows to save the given reaction into the database.
func (db Db) SavePostReaction(reaction types.PostReaction) error {
	err := db.SaveUserIfNotExisting(reaction.Owner, reaction.Height)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO post_reaction (post_id, owner_address, short_code, value, height) 
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT ON CONSTRAINT react_unique DO UPDATE 
    SET post_id = excluded.post_id, 
        owner_address = excluded.owner_address,
        short_code = excluded.short_code, 
		value = excluded.value,
		height = excluded.height
WHERE post_reaction.height <= excluded.height`
	_, err = db.Sql.Exec(stmt, reaction.PostID, reaction.Owner, reaction.ShortCode, reaction.Value, reaction.Height)
	return err
}

// RemovePostReaction allows to remove an already existing reaction from the database.
func (db Db) RemovePostReaction(reaction types.PostReaction) error {
	err := db.SaveUserIfNotExisting(reaction.Owner, reaction.Height)
	if err != nil {
		return err
	}

	statement := `
DELETE FROM post_reaction 
WHERE post_id = $1 AND owner_address = $2 AND short_code = $3 AND height <= $4`
	_, err = db.Sql.Exec(statement, reaction.PostID, reaction.Owner, reaction.ShortCode, reaction.Height)
	return err
}

// GetRegisteredReactionByCodeOrValue allows to get a registered reaction by its shortcode or
// value and the subspace for which it has been registered.
func (db Db) GetRegisteredReactionByCodeOrValue(
	codeOrValue string, subspace string,
) (*poststypes.RegisteredReaction, error) {
	stmt := `SELECT * FROM registered_reactions WHERE (short_code = $1 OR value = $1) AND subspace = $2`

	var rows []dbtypes.RegisteredReactionRow
	err := db.Sqlx.Select(&rows, stmt, codeOrValue, subspace)
	if err != nil {
		return nil, err
	}

	// No post found
	if len(rows) == 0 {
		return nil, nil
	}

	reaction := convertReactionRow(rows[0])
	return &reaction, nil
}

// RegisterReactionIfNotPresent allows to register into the database the given reaction.
func (db Db) RegisterReactionIfNotPresent(reaction types.RegisteredReaction) error {
	react, err := db.GetRegisteredReactionByCodeOrValue(reaction.ShortCode, reaction.Subspace)
	if err != nil {
		return err
	}

	// If the reaction exists do nothing
	if react != nil {
		return nil
	}

	// Save the owner
	err = db.SaveUserIfNotExisting(reaction.Creator, reaction.Height)
	if err != nil {
		return err
	}

	// Save the reaction
	stmt := `
INSERT INTO registered_reactions (owner_address, short_code, value, subspace, height) 
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT ON CONSTRAINT registered_react_unique DO UPDATE 
    SET owner_address = excluded.owner_address,
    	short_code = excluded.short_code,
    	value = excluded.value, 
    	subspace = excluded.subspace,
    	height = excluded.height
WHERE registered_reactions.height <= excluded.height`
	_, err = db.Sql.Exec(stmt, reaction.Creator, reaction.ShortCode, reaction.Value, reaction.Subspace, reaction.Height)
	return err
}
