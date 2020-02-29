package handlers

import (
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/juno/db/postgresql"
)

func handleMsgRemovePostReaction(postID uint64, msg posts.MsgRemovePostReaction, db postgresql.Database) error {

	removeRSqlStatement := `
	DELETE FROM reaction
	WHERE id = $1 AND owner = $2 AND val = $3;
	`
	_, err := db.Sql.Exec(
		removeRSqlStatement,
		msg.PostID,
		msg.User,
		msg.Reaction,
	)

	if err != nil {
		return err
	}

	return nil
}
