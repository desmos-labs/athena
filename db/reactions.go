package db

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	emoji "github.com/desmos-labs/Go-Emoji-Utils"
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/rs/zerolog/log"
)

// PostRow represents a single PostgreSQL row containing the data of a Post
type RegisteredReactionRow struct {
	ReactionID uint64  `db:"id"`
	OwnerID    *uint64 `db:"owner_id"`
	ShortCode  string  `db:"short_code"`
	Value      string  `db:"value"`
	Subspace   string  `db:"subspace"`
}

// ConvertPostRow takes the given postRow and userRow and merges the data contained inside them to create a Post.
func ConvertReactionRow(reactionRow RegisteredReactionRow, userRow *UserRow) (*posts.Reaction, error) {

	// Parse the creator
	creator, err := sdk.AccAddressFromBech32(userRow.Address)
	if err != nil {
		return nil, err
	}

	// Create the reaction
	reaction := posts.NewReaction(creator, reactionRow.ShortCode, reactionRow.Value, reactionRow.Subspace)
	return &reaction, nil
}

// SaveReaction allows to save a new reaction for the given postID having the specified value and user
func (db DesmosDb) SaveReaction(postID posts.PostID, reaction posts.PostReaction) (*posts.PostReaction, error) {
	owner, err := db.SaveUserIfNotExisting(reaction.Owner)
	if err != nil {
		return nil, err
	}

	// Get the post
	post, err := db.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	// Get the registered reaction from the reaction object
	var registeredReact *posts.Reaction
	if e, err := emoji.LookupEmoji(reaction.Value); err == nil {
		// If the reaction is an emoji, register the reaction
		reaction := posts.NewReaction(posts.ModuleAddress, e.Shortcodes[0], e.Value, post.Subspace)
		registeredReact, err = db.RegisterReactionIfNotPresent(reaction)
		if err != nil {
			return nil, err
		}
	} else {
		// The reaction is not an emoji, so read it from the database
		registeredReact, err = db.GetRegisteredReactionByCodeOrValue(reaction.Value, post.Subspace)
		if err != nil {
			return nil, err
		}
	}

	log.Info().
		Str("post_id", postID.String()).
		Str("value", registeredReact.Value).
		Str("short_code", registeredReact.ShortCode).
		Msg("saving reaction")

	statement := `INSERT INTO reaction (post_id, owner_id, short_code, value) VALUES ($1, $2, $3, $4)`
	_, err = db.Sql.Exec(statement, postID, owner.Id, registeredReact.ShortCode, registeredReact.Value)
	return &reaction, err
}

// RemoveReaction allows to remove an already existing reaction for the post having the given postID,
// the given reaction and from the specified user.
func (db DesmosDb) RemoveReaction(postID posts.PostID, reaction posts.PostReaction) error {
	owner, err := db.SaveUserIfNotExisting(reaction.Owner)
	if err != nil {
		return err
	}

	log.Info().
		Str("post_id", postID.String()).
		Str("value", reaction.Value).
		Msg("removing reaction")

	statement := `DELETE FROM reaction WHERE post_id = $1 AND owner_id = $2 AND (value = $3 OR short_code = $3)`
	_, err = db.Sql.Exec(statement, postID.String(), owner.Id, reaction.Value)
	return err
}

// GetRegisteredReactionByCodeOrValue allows to get a registered reaction by its shortcode or
// value and the subspace for which it has been registered.
func (db DesmosDb) GetRegisteredReactionByCodeOrValue(codeOrValue string, subspace string) (*posts.Reaction, error) {
	postSqlStatement := `SELECT * FROM registered_reactions WHERE (short_code = $1 OR value = $1) AND subspace = $2`

	var rows []RegisteredReactionRow
	err := db.sqlx.Select(&rows, postSqlStatement, codeOrValue, subspace)
	if err != nil {
		return nil, err
	}

	// No post found
	if len(rows) == 0 {
		return nil, nil
	}

	reactionRow := rows[0]

	// Find the user
	userRow, err := db.GetUserById(reactionRow.OwnerID)
	if err != nil {
		return nil, err
	}

	return ConvertReactionRow(reactionRow, userRow)
}

// RegisterReaction allows to register into the database the given reaction.
func (db DesmosDb) RegisterReactionIfNotPresent(reaction posts.Reaction) (*posts.Reaction, error) {
	react, err := db.GetRegisteredReactionByCodeOrValue(reaction.ShortCode, reaction.Subspace)
	if err != nil {
		return nil, err
	}

	// If the reaction exists do nothing
	if react != nil {
		return react, nil
	}

	// Save the owner
	owner, err := db.SaveUserIfNotExisting(reaction.Creator)
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("value", reaction.Value).
		Str("short_code", reaction.ShortCode).
		Str("creator", reaction.Creator.String()).
		Msg("registering reaction")

	// Save the reaction
	statement := `INSERT INTO registered_reactions (owner_id, short_code, value, subspace) VALUES ($1, $2, $3, $4)`
	_, err = db.Sql.Exec(statement, owner.Id, reaction.ShortCode, reaction.Value, reaction.Subspace)
	return &reaction, err
}
