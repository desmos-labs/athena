package notifications

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	poststypes "github.com/desmos-labs/desmos/v4/x/posts/types"
	reactionstypes "github.com/desmos-labs/desmos/v4/x/reactions/types"
	relationshipstypes "github.com/desmos-labs/desmos/v4/x/relationships/types"
	juno "github.com/forbole/juno/v3/types"
)

// HandleMsg implements modules.MessageModule
func (m Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch desmosMsg := msg.(type) {
	case *relationshipstypes.MsgCreateRelationship:
		return m.handleMsgCreateRelationship(desmosMsg)

	case *poststypes.MsgCreatePost:
		return m.handleMsgCreatePost(tx, index, desmosMsg)

	case *reactionstypes.MsgAddReaction:
		return m.handleMsgAddReaction(tx, desmosMsg)
	}

	return nil
}

// handleMsgCreateRelationship handles a MsgCreateRelationship message and sends out the various related notifications
func (m Module) handleMsgCreateRelationship(msg *relationshipstypes.MsgCreateRelationship) error {
	// Skip if the subspace is not the correct one
	if msg.SubspaceID != m.cfg.SubspaceID {
		return nil
	}

	return m.SendRelationshipNotifications(msg.SubspaceID, msg.Signer, msg.Counterparty)
}

// handleMsgCreatePost handles a MsgCreatePost message and sends out the various related notifications
func (m Module) handleMsgCreatePost(tx *juno.Tx, index int, msg *poststypes.MsgCreatePost) error {
	// Skip if the subspace is not the correct one
	if msg.SubspaceID != m.cfg.SubspaceID {
		return nil
	}

	// Get the post id
	event, err := tx.FindEventByType(index, poststypes.EventTypeCreatePost)
	if err != nil {
		return err
	}
	postIDStr, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyPostID)
	if err != nil {
		return err
	}
	postID, err := poststypes.ParsePostID(postIDStr)
	if err != nil {
		return err
	}

	// Send the notifications
	return m.SendPostNotifications(tx.Height, msg.SubspaceID, postID)
}

// handleMsgAddReaction handles a MsgAddReaction message and sends out the various related notifications
func (m Module) handleMsgAddReaction(tx *juno.Tx, msg *reactionstypes.MsgAddReaction) error {
	// Skip if the subspace is not the correct one
	if msg.SubspaceID != m.cfg.SubspaceID {
		return nil
	}

	// Send the notifications
	return m.SendReactionNotifications(tx.Height, msg.SubspaceID, msg.PostID, msg.User)
}
