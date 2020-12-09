package notifications

import (
	poststypes "github.com/desmos-labs/desmos/x/posts/types"
	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"
	"github.com/desmos-labs/djuno/notifications"
	juno "github.com/desmos-labs/juno/types"
)

// TxHandler allows to handle a transaction in order to send the
func TxHandler(tx *juno.Tx) error {
	if hasDesmosMsg, desmosUser := getDesmosUser(tx); hasDesmosMsg {
		return notifications.SendTransactionResultNotification(tx, desmosUser)
	}
	return nil
}

// getDesmosUser returns the first Desmos address that has created a Desmos message
// inside the given transaction. If no Desmos message could be found, returns false.
func getDesmosUser(tx *juno.Tx) (bool, string) {
	// TODO: Add other message types
	for _, msg := range tx.Msgs {
		switch desmosMsg := msg.(type) {
		// Posts
		case *poststypes.MsgCreatePost:
			return true, desmosMsg.Creator
		case *poststypes.MsgEditPost:
			return true, desmosMsg.Editor

		// Reactions
		case *poststypes.MsgRegisterReaction:
			return true, desmosMsg.Creator
		case *poststypes.MsgAddPostReaction:
			return true, desmosMsg.User
		case *poststypes.MsgRemovePostReaction:
			return true, desmosMsg.User

		// Polls
		case *poststypes.MsgAnswerPoll:
			return true, desmosMsg.Answerer

		// Profiles
		case *profilestypes.MsgSaveProfile:
			return true, desmosMsg.Creator
		}
	}
	return false, ""
}
