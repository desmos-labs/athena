package db

import (
	"github.com/desmos-labs/desmos/x/posts"
)

// SaveReaction allows to save a new reaction for the given postID having the specified value and user
func (db DesmosDb) SaveReaction(postID posts.PostID, reaction posts.Reaction) (*posts.Reaction, error) {
	owner, err := db.SaveUserIfNotExisting(reaction.Owner)
	if err != nil {
		return nil, err
	}

	statement := `INSERT INTO reaction (post_id, owner_id, value) VALUES ($1, $2, $3)`
	_, err = db.Sql.Exec(statement, postID, owner.Id, reaction.Value)
	return &reaction, err
}

// RemoveReaction allows to remove an already existing reaction for the post having the given postID,
// the given reaction and from the specified user.
func (db DesmosDb) RemoveReaction(postID posts.PostID, reaction posts.Reaction) error {
	owner, err := db.SaveUserIfNotExisting(reaction.Owner)
	if err != nil {
		return err
	}

	statement := `DELETE FROM reaction WHERE post_id = $1 AND owner_id = $2 AND reaction = $3;`
	_, err = db.Sql.Exec(statement, postID, owner.Id, reaction)
	return err
}
