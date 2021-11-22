package notifications

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	poststypes "github.com/desmos-labs/desmos/v2/x/staging/posts/types"
	juno "github.com/forbole/juno/v2/types"

	postsutils "github.com/desmos-labs/djuno/v2/x/posts/utils"
)

// HandleMsg implements modules.MessageModule
func (m Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch desmosMsg := msg.(type) {
	// Posts
	case *poststypes.MsgCreatePost:
		return m.sendPostNotification(tx, index, desmosMsg)

	case *poststypes.MsgAddPostReaction:
		return m.sendReactionNotification(tx, index, poststypes.EventTypePostReactionAdded)

	case *poststypes.MsgRemovePostReaction:
		return m.sendReactionNotification(tx, index, poststypes.EventTypePostReactionRemoved)
	}

	return nil
}

// sendPostNotification sends a notification to everyone that might be involved in a post (eg. tags, etc)
func (m *Module) sendPostNotification(tx *juno.Tx, index int, msg *poststypes.MsgCreatePost) error {
	post, err := postsutils.GetPostFromMsgCreatePost(tx, index, msg)
	if err != nil {
		return err
	}

	return m.sendPostNotifications(post)
}

// sendReactionNotification sends a notification to the creator of the post to which the reaction has been added
func (m *Module) sendReactionNotification(tx *juno.Tx, index int, event string) error {
	postID, reaction, err := postsutils.GetReactionFromTxEvent(tx, index, event)
	if err != nil {
		return err
	}

	return m.sendReactionNotifications(postID, reaction)
}
