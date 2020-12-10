package database

import (
	poststypes "github.com/desmos-labs/desmos/x/posts/types"
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
func (db DesmosDb) SavePostReaction(postID string, reaction *poststypes.PostReaction) error {
	err := db.SaveUserIfNotExisting(reaction.Owner)
	if err != nil {
		return err
	}

	statement := `INSERT INTO reaction (post_id, owner_address, short_code, value) VALUES ($1, $2, $3, $4)`
	_, err = db.Sql.Exec(statement, postID, reaction.Owner, reaction.ShortCode, reaction.Value)
	return err
}

// RemoveReaction allows to remove an already existing reaction from the database.
func (db DesmosDb) RemoveReaction(postID string, reaction *poststypes.PostReaction) error {
	err := db.SaveUserIfNotExisting(reaction.Owner)
	if err != nil {
		return err
	}

	statement := `DELETE FROM reaction WHERE post_id = $1 AND owner_address = $2 AND short_code = $3`
	_, err = db.Sql.Exec(statement, postID, reaction.Owner, reaction.ShortCode)
	return err
}

// GetRegisteredReactionByCodeOrValue allows to get a registered reaction by its shortcode or
// value and the subspace for which it has been registered.
func (db DesmosDb) GetRegisteredReactionByCodeOrValue(
	codeOrValue string, subspace string,
) (*poststypes.RegisteredReaction, error) {
	postSqlStatement := `SELECT * FROM registered_reactions WHERE (short_code = $1 OR value = $1) AND subspace = $2`

	var rows []dbtypes.RegisteredReactionRow
	err := db.Sqlx.Select(&rows, postSqlStatement, codeOrValue, subspace)
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

// RegisterReaction allows to register into the database the given reaction.
func (db DesmosDb) RegisterReactionIfNotPresent(reaction poststypes.RegisteredReaction) error {
	react, err := db.GetRegisteredReactionByCodeOrValue(reaction.ShortCode, reaction.Subspace)
	if err != nil {
		return err
	}

	// If the reaction exists do nothing
	if react != nil {
		return nil
	}

	// Save the owner
	err = db.SaveUserIfNotExisting(reaction.Creator)
	if err != nil {
		return err
	}

	// Save the reaction
	statement := `INSERT INTO registered_reactions (owner_address, short_code, value, subspace) VALUES ($1, $2, $3, $4)`
	_, err = db.Sql.Exec(statement, reaction.Creator, reaction.ShortCode, reaction.Value, reaction.Subspace)
	return err
}
