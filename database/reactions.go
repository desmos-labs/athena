package database

import (
	"fmt"

	"github.com/desmos-labs/djuno/v2/types"
)

// SaveReaction stores the given reaction inside the database
func (db *Db) SaveReaction(reaction types.Reaction) error {
	stmt := `
INSERT INTO reaction (subspace_id, post_id, id, value, author_address, height) 
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT DO UPDATE 
    SET value = excluded.value,
        author_address = excluded.author_address,
        height = excluded.height
WHERE reaction.height <= excluded.height`

	valueBz, err := db.EncodingConfig.Marshaler.MarshalJSON(reaction.Value)
	if err != nil {
		return fmt.Errorf("failed to json encode reaction value: %s", err)
	}

	_, err = db.Sql.Exec(stmt,
		reaction.SubspaceID,
		reaction.PostID,
		reaction.ID,
		string(valueBz),
		reaction.Author,
		reaction.Height,
	)
	return err
}

// DeleteReaction removes the given reaction from the database
func (db *Db) DeleteReaction(height int64, subspaceID uint64, postID uint64, reactionID uint32) error {
	stmt := `DELETE FROM reaction WHERE subspace_id = $1 AND post_id = $2 AND id = $3 AND height <= $4`
	_, err := db.Sql.Exec(stmt, subspaceID, postID, reactionID, height)
	return err
}

// --------------------------------------------------------------------------------------------------------------------

// SaveRegisteredReaction stores the given registered reaction inside the database
func (db *Db) SaveRegisteredReaction(reaction types.RegisteredReaction) error {
	stmt := `
INSERT INTO registered_reaction (subspace_id, id, shorthand_code, display_value, height) 
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT DO UPDATE 
    SET shorthand_code = excluded.shorthand_code,
        display_value = excluded.display_value,
        height = excluded.height
WHERE registered_reaction.height <= excluded.height`

	_, err := db.Sql.Exec(stmt,
		reaction.SubspaceID,
		reaction.ID,
		reaction.ShorthandCode,
		reaction.DisplayValue,
		reaction.Height,
	)
	return err
}

// DeleteRegisteredReaction removes the given registered reaction from the database
func (db *Db) DeleteRegisteredReaction(height int64, subspaceID uint64, reactionID uint32) error {
	stmt := `DELETE FROM registered_reaction WHERE subspace_id = $1 AND id = $2 AND height <= $3`
	_, err := db.Sql.Exec(stmt, subspaceID, reactionID, height)
	return err
}

// --------------------------------------------------------------------------------------------------------------------

// SaveReactionParams saves the given params inside the database
func (db *Db) SaveReactionParams(params types.ReactionParams) error {
	// Store registered reaction params
	stmt := `
INSERT INTO subspace_registered_reaction_params (subspace_id, enabled, height) 
VALUES ($1, $2, $3)
ON CONFLICT DO UPDATE 
    SET enabled = excluded.enabled,
        height = excluded.height
WHERE subspace_registered_reaction_params.height <= excluded.height`

	_, err := db.Sql.Exec(stmt, params.SubspaceID, params.RegisteredReaction.Enabled, params.Height)
	if err != nil {
		return err
	}

	// Store free text params
	stmt = `
INSERT INTO subspace_free_text_params (subspace_id, enabled, max_length, reg_ex, height) 
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT DO UPDATE 
    SET enabled = excluded.enabled,
        max_length = excluded.max_length,
        reg_ex = excluded.reg_ex,
        height = excluded.height
WHERE subspace_free_text_params.height <= excluded.height`

	_, err = db.Sql.Exec(stmt,
		params.SubspaceID,
		params.FreeText.Enabled,
		params.FreeText.MaxLength,
		params.FreeText.RegEx,
		params.Height,
	)
	return err
}
