package handlers

import (
	poststypes "github.com/desmos-labs/desmos/x/posts/types"
	"github.com/desmos-labs/djuno/database"
)

// HandleMsgAnswerPoll allows to properly handle a MsgAnswerPoll message by
// storing inside the database the new answer.
func HandleMsgAnswerPoll(msg poststypes.MsgAnswerPoll, db database.DesmosDb) error {
	return db.SaveUserPollAnswer(msg.PostID, poststypes.NewUserAnswer(msg.UserAnswers, msg.Answerer))
}
