package handlers

import (
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/juno/db/postgresql"
)

func handleMsgRemovePostReaction(msg posts.MsgRemovePostReaction, db postgresql.Database) error {

	removeRSqlStatement := `
	DELETE FROM reaction
	WHERE post_id = $1 AND owner = $2 AND value = $3;
	`
	_, err := db.Sql.Exec(
		removeRSqlStatement,
		msg.PostID,
		msg.User.String(),
		msg.Reaction,
	)

	if err != nil {
		return err
	}

	return nil
}
