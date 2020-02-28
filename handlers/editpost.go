package handlers

import (
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/juno/db/postgresql"
)

func handleMsgEditPost(postID uint64, msg posts.MsgEditPost, db postgresql.Database) error {
	var id uint64

	// Editing Post

	postSqlStatement := `
	UPDATE post 
	SET message = $1, last_edited = $2
	WHERE id = $3 AND creator = $4
	RETURNING id;
	`

	err := db.Sql.QueryRow(
		postSqlStatement,
		msg.Message,
		msg.EditDate,
		msg.PostID,
		msg.Editor,
	).Scan(&id)

	if err != nil {
		return err
	}

	return nil
}
