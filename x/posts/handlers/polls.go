package handlers

import (
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/djuno/database"
)

// HandleMsgAnswerPoll allows to properly handle a MsgAnswerPoll message by
// storing inside the database the new answer.
func HandleMsgAnswerPoll(msg posts.MsgAnswerPoll, db database.DesmosDb) error {
	return db.SavePollAnswer(msg.PostID, posts.NewUserAnswer(msg.UserAnswers, msg.Answerer))
}
