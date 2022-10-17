package reactions

import (
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/gogo/protobuf/proto"

	reactionstypes "github.com/desmos-labs/desmos/v4/x/reactions/types"

	"github.com/rs/zerolog/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v3/types"
)

// HandleMsgExec implements modules.AuthzMessageModule
func (m *Module) HandleMsgExec(index int, _ *authz.MsgExec, _ int, executedMsg sdk.Msg, tx *juno.Tx) error {
	return m.HandleMsg(index, executedMsg, tx)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch desmosMsg := msg.(type) {
	case *reactionstypes.MsgAddReaction:
		return m.handleMsgAddReaction(tx, index, desmosMsg)

	case *reactionstypes.MsgRemoveReaction:
		return m.handleMsgRemoveReaction(tx, desmosMsg)

	case *reactionstypes.MsgAddRegisteredReaction:
		return m.handleMsgAddRegisteredReaction(tx, index, desmosMsg)

	case *reactionstypes.MsgEditRegisteredReaction:
		return m.handleMsgEditRegisteredReaction(tx, desmosMsg)

	case *reactionstypes.MsgRemoveRegisteredReaction:
		return m.handleMsgRemoveRegisteredReaction(tx, desmosMsg)

	case *reactionstypes.MsgSetReactionsParams:
		return m.handleMsgSetReactionsParams(tx, desmosMsg)
	}

	log.Debug().Str("module", "reactions").Str("message", proto.MessageName(msg)).
		Int64("height", tx.Height).Msg("handled message")

	return nil
}

// handleMsgAddReaction handles a MsgAddReaction
func (m *Module) handleMsgAddReaction(tx *juno.Tx, index int, msg *reactionstypes.MsgAddReaction) error {
	reactionID, err := m.GetReactionID(tx, index)
	if err != nil {
		return err
	}

	reaction, err := m.GetReaction(tx.Height, msg.SubspaceID, msg.PostID, reactionID)
	if err != nil {
		return err
	}

	return m.db.SaveReaction(reaction)
}

// handleMsgRemoveReaction handles a MsgRemoveReaction
func (m *Module) handleMsgRemoveReaction(tx *juno.Tx, msg *reactionstypes.MsgRemoveReaction) error {
	return m.db.DeleteReaction(tx.Height, msg.SubspaceID, msg.PostID, msg.ReactionID)
}

// handleMsgAddRegisteredReaction handles a MsgAddRegisteredReaction
func (m *Module) handleMsgAddRegisteredReaction(tx *juno.Tx, index int, msg *reactionstypes.MsgAddRegisteredReaction) error {
	event, err := tx.FindEventByType(index, reactionstypes.EventTypeAddRegisteredReaction)
	if err != nil {
		return err
	}
	reactionIDStr, err := tx.FindAttributeByKey(event, reactionstypes.AttributeKeyRegisteredReactionID)
	if err != nil {
		return err
	}
	reactionID, err := reactionstypes.ParseRegisteredReactionID(reactionIDStr)
	if err != nil {
		return err
	}

	return m.updateRegisteredReaction(tx.Height, msg.SubspaceID, reactionID)
}

// handleMsgEditRegisteredReaction handles a MsgEditRegisteredReaction
func (m *Module) handleMsgEditRegisteredReaction(tx *juno.Tx, msg *reactionstypes.MsgEditRegisteredReaction) error {
	return m.updateRegisteredReaction(tx.Height, msg.SubspaceID, msg.RegisteredReactionID)
}

// handleMsgRemoveRegisteredReaction handles a MsgRemoveRegisteredReacton
func (m *Module) handleMsgRemoveRegisteredReaction(tx *juno.Tx, msg *reactionstypes.MsgRemoveRegisteredReaction) error {
	return m.db.DeleteRegisteredReaction(tx.Height, msg.SubspaceID, msg.RegisteredReactionID)
}

// handleMsgSetReactionsParams handles a MsgSetReactionsParams
func (m *Module) handleMsgSetReactionsParams(tx *juno.Tx, msg *reactionstypes.MsgSetReactionsParams) error {
	return m.updateReactionParams(tx.Height, msg.SubspaceID)
}
