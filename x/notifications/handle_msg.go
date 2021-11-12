package notifications

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	poststypes "github.com/desmos-labs/desmos/v2/x/staging/posts/types"
	juno "github.com/forbole/juno/v2/types"

	"github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/djuno/x/notifications/utils"
	postsutils "github.com/desmos-labs/djuno/x/posts/utils"
)

// MsgHandler allows to handle different message types from the posts module
func MsgHandler(tx *juno.Tx, index int, msg sdk.Msg, database *database.Db) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch desmosMsg := msg.(type) {
	// Posts
	case *poststypes.MsgCreatePost:
		return sendPostNotification(tx, index, desmosMsg, database)

	case *poststypes.MsgAddPostReaction:
		return sendReactionNotification(tx, index, poststypes.EventTypePostReactionAdded, database)

	case *poststypes.MsgRemovePostReaction:
		return sendReactionNotification(tx, index, poststypes.EventTypePostReactionRemoved, database)
	}

	return nil
}

// sendPostNotification sends a notification to everyone that might be involved in a post (eg. tags, etc)
func sendPostNotification(tx *juno.Tx, index int, msg *poststypes.MsgCreatePost, db *database.Db) error {
	post, err := postsutils.GetPostFromMsgCreatePost(tx, index, msg)
	if err != nil {
		return err
	}

	return utils.SendPostNotifications(post, db)
}

// sendReactionNotification sends a notification to the creator of the post to which the reaction has been added
func sendReactionNotification(tx *juno.Tx, index int, event string, db *database.Db) error {
	postID, reaction, err := postsutils.GetReactionFromTxEvent(tx, index, event)
	if err != nil {
		return err
	}

	return utils.SendReactionNotifications(postID, reaction, db)
}
