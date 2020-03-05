package handlers

import (
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/juno/db/postgresql"
)

func handleMsgRemovePostReaction(msg posts.MsgRemovePostReaction, db postgresql.Database) error {
	var id uint64

	removeRSqlStatement := `
	DELETE FROM reaction
	WHERE post_id = $1 AND owner = $2 AND value = $3;
	`
	return db.Sql.QueryRow(
		removeRSqlStatement,
		msg.PostID,
		msg.User.String(),
		msg.Reaction,
	).Scan(&id)
}
