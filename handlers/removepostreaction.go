package handlers

import (
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/juno/db/postgresql"
)

func handleMsgRemovePostReaction(postID uint64, msg posts.MsgRemovePostReaction, db postgresql.Database) error {
	return nil
}
