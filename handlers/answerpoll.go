package handlers

import (
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/juno/db/postgresql"
)

func handleMsgAnswerPoll(postID uint64, msg posts.MsgAnswerPoll, db postgresql.Database) error {
	return nil
}
