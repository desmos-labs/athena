package handlers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/djuno/notifications"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/types"
)

// TxHandler handles each single transaction, verifying if it has been successful.
// If not, notifies the user that has created it using the notifications system.
func TxHandler(tx types.Tx, _ db.Database) error {
	if hasDesmosMsg, desmosUser := getDesmosUser(tx); hasDesmosMsg {
		return notifications.SendTransactionResultNotification(tx, desmosUser)
	}
	return nil
}

// getDesmosUser returns the first Desmos address that has created a Desmos message
// inside the given transaction. If no Desmos message could be found, returns false.
func getDesmosUser(tx types.Tx) (bool, sdk.AccAddress) {
	for _, msg := range tx.Messages {
		switch desmosMsg := msg.(type) {
		case posts.MsgCreatePost:
			return true, desmosMsg.Creator
		case posts.MsgEditPost:
			return true, desmosMsg.Editor
		case posts.MsgAddPostReaction:
			return true, desmosMsg.User
		case posts.MsgRemovePostReaction:
			return true, desmosMsg.User
		case posts.MsgAnswerPoll:
			return true, desmosMsg.Answerer
		}
	}
	return false, nil
}
