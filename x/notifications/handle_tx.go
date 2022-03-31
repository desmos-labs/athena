package notifications

import (
	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"
	poststypes "github.com/desmos-labs/desmos/v2/x/staging/posts/types"
	juno "github.com/forbole/juno/v3/types"
)

// HandleTx implements modules.TransactionModule
func (m *Module) HandleTx(tx *juno.Tx) error {
	if hasDesmosMsg, desmosUser := getDesmosUser(tx); hasDesmosMsg {
		return m.sendTransactionResultNotification(tx, desmosUser)
	}
	return nil
}

// getDesmosUser returns the first Desmos address that has created a Desmos message
// inside the given transaction. If no Desmos message could be found, returns false.
func getDesmosUser(tx *juno.Tx) (bool, string) {
	// TODO: Add other message types
	for _, msg := range tx.GetMsgs() {
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
