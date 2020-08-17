package database

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	poststypes "github.com/desmos-labs/desmos/x/posts/types"
	dbtypes "github.com/desmos-labs/djuno/database/types"

	"github.com/rs/zerolog/log"
)

// convertPostRow takes the given postRow and userRow and merges the data contained inside them to create a Post.
func convertReactionRow(reactionRow dbtypes.RegisteredReactionRow, userRow *dbtypes.ProfileRow) (*poststypes.Reaction, error) {

	// Parse the creator
	creator, err := sdk.AccAddressFromBech32(userRow.Address)
	if err != nil {
		return nil, err
	}

	// Create the reaction
	reaction := poststypes.NewReaction(creator, reactionRow.ShortCode, reactionRow.Value, reactionRow.Subspace)
	return &reaction, nil
}

// GetRegisteredReactionByCodeOrValue allows to get a registered reaction by its shortcode or
// value and the subspace for which it has been registered.
func (db DesmosDb) GetRegisteredReactionByCodeOrValue(codeOrValue string, subspace string) (*poststypes.Reaction, error) {
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

	reactionRow := rows[0]

	// Find the user
	addr, err := sdk.AccAddressFromBech32(reactionRow.OwnerAddress)
	if err != nil {
		return nil, err
	}

	userRow, err := db.GetUserByAddress(addr)
	if err != nil {
		return nil, err
	}

	return convertReactionRow(reactionRow, userRow)
}

// ______________________________________________

// SaveReaction allows to save the given reaction into the database.
func (db DesmosDb) SaveReaction(postID poststypes.PostID, reaction *poststypes.PostReaction) error {
	if err := db.SaveUserIfNotExisting(reaction.Owner); err != nil {
		return err
	}

	log.Info().
		Str("module", "poststypes").
		Str("post_id", postID.String()).
		Str("value", reaction.Value).
		Str("short_code", reaction.Shortcode).
		Str("user", reaction.Owner.String()).
		Msg("saving reaction")

	stmt := `INSERT INTO reaction (post_id, owner_address, short_code, value) 
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT ON CONSTRAINT react_unique DO UPDATE 
			     SET post_id = excluded.post_id, 
			         owner_address = excluded.owner_address,
			         short_code = excluded.short_code,
			         value = excluded.value`
	_, err := db.Sql.Exec(stmt, postID.String(), reaction.Owner.String(), reaction.Shortcode, reaction.Value)
	return err
}

// RegisterReaction allows to register into the database the given reaction.
func (db DesmosDb) SaveRegisteredReactionIfNotPresent(reaction poststypes.Reaction) (*poststypes.Reaction, error) {
	react, err := db.GetRegisteredReactionByCodeOrValue(reaction.ShortCode, reaction.Subspace)
	if err != nil {
		return nil, err
	}

	// If the reaction exists do nothing
	if react != nil {
		return react, nil
	}

	// Save the owner
	if err := db.SaveUserIfNotExisting(reaction.Creator); err != nil {
		return nil, err
	}

	log.Info().
		Str("value", reaction.Value).
		Str("short_code", reaction.ShortCode).
		Str("creator", reaction.Creator.String()).
		Msg("registering reaction")

	// Save the reaction
	statement := `INSERT INTO registered_reactions (owner_address, short_code, value, subspace) VALUES ($1, $2, $3, $4)`
	_, err = db.Sql.Exec(statement, reaction.Creator.String(), reaction.ShortCode, reaction.Value, reaction.Subspace)
	return &reaction, err
}

// ______________________________________________

// RemoveReaction allows to remove an already existing reaction from the database.
func (db DesmosDb) RemoveReaction(postID poststypes.PostID, reaction *poststypes.PostReaction) error {
	if err := db.SaveUserIfNotExisting(reaction.Owner); err != nil {
		return err
	}

	log.Info().
		Str("post_id", postID.String()).
		Str("value", reaction.Value).
		Str("short_code", reaction.Shortcode).
		Str("user", reaction.Owner.String()).
		Msg("removing reaction")

	statement := `DELETE FROM reaction WHERE post_id = $1 AND owner_address = $2 AND short_code = $3`
	_, err := db.Sql.Exec(statement, postID.String(), reaction.Owner.String(), reaction.Shortcode)
	return err
}
