package handlers

import (
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/juno/db/postgresql"
)

func handleMsgAddPostReaction(msg posts.MsgAddPostReaction, db postgresql.Database) error {
	var id uint64

	addReactionSqlStatement := `
	INSERT INTO reaction (post_id, owner, value)
	VALUES ($1, $2, $3)
	RETURNING id;
	`

	return db.Sql.QueryRow(
		addReactionSqlStatement,
		msg.PostID,
		msg.User.String(),
		msg.Value,
	).Scan(&id)

}
