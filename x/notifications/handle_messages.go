package notifications

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	poststypes "github.com/desmos-labs/desmos/v6/x/posts/types"
	reactionstypes "github.com/desmos-labs/desmos/v6/x/reactions/types"
	relationshipstypes "github.com/desmos-labs/desmos/v6/x/relationships/types"
	juno "github.com/forbole/juno/v5/types"

	"github.com/desmos-labs/athena/v2/x/filters"

	"github.com/desmos-labs/athena/v2/types"
)

// HandleMsgExec implements modules.AuthzMessageModule
func (m *Module) HandleMsgExec(index int, _ *authz.MsgExec, _ int, executedMsg sdk.Msg, tx *juno.Tx) error {
	return m.HandleMsg(index, executedMsg, tx)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 || !filters.ShouldMsgBeParsed(msg) {
		return nil
	}

	switch desmosMsg := msg.(type) {
	case *relationshipstypes.MsgCreateRelationship:
		return m.handleMsgCreateRelationship(tx, desmosMsg)

	case *poststypes.MsgCreatePost:
		return m.handleMsgCreatePost(tx, index, desmosMsg)

	case *reactionstypes.MsgAddReaction:
		return m.handleMsgAddReaction(tx, index, desmosMsg)
	}

	return nil
}

// handleMsgCreateRelationship handles a MsgCreateRelationship message and sends out the various related notifications
func (m *Module) handleMsgCreateRelationship(tx *juno.Tx, msg *relationshipstypes.MsgCreateRelationship) error {
	return m.SendRelationshipNotifications(types.NewRelationship(
		relationshipstypes.NewRelationship(msg.Signer, msg.Counterparty, msg.SubspaceID),
		tx.Height,
	))
}

// handleMsgCreatePost handles a MsgCreatePost message and sends out the various related notifications
func (m *Module) handleMsgCreatePost(tx *juno.Tx, index int, msg *poststypes.MsgCreatePost) error {
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
func (m *Module) handleMsgAddReaction(tx *juno.Tx, index int, msg *reactionstypes.MsgAddReaction) error {
	// Get the reaction value
	reactionID, err := m.reactionsModule.GetReactionID(tx, index)
	if err != nil {
		return err
	}
	reaction, err := m.reactionsModule.GetReaction(tx.Height, msg.SubspaceID, msg.PostID, reactionID)
	if err != nil {
		return err
	}

	// Send the notifications
	return m.SendReactionNotifications(reaction)
}
