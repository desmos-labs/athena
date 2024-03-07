package relationships

import (
	"github.com/desmos-labs/athena/v2/x/filters"

	"github.com/cosmos/cosmos-sdk/x/authz"

	relationshipstypes "github.com/desmos-labs/desmos/v7/x/relationships/types"

	"github.com/rs/zerolog/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v5/types"

	"github.com/desmos-labs/athena/v2/types"
)

// HandleMsgExec implements modules.AuthzMessageModule
func (m *Module) HandleMsgExec(index int, _ *authz.MsgExec, _ int, executedMsg sdk.Msg, tx *juno.Tx) error {
	return m.HandleMsg(index, executedMsg, tx)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 || !filters.ShouldMsgBeParsed(msg) {
		return nil
	}

	switch desmosMsg := msg.(type) {

	case *relationshipstypes.MsgCreateRelationship:
		return m.handleMsgCreateRelationship(tx, desmosMsg)

	case *relationshipstypes.MsgDeleteRelationship:
		return m.handleMsgDeleteRelationship(tx, desmosMsg)

	case *relationshipstypes.MsgBlockUser:
		return m.handleMsgBlockUser(tx, desmosMsg)

	case *relationshipstypes.MsgUnblockUser:
		return m.handleMsgUnblockUser(tx, desmosMsg)

	}

	log.Debug().Str("module", "relationships").Str("message", sdk.MsgTypeURL(msg)).
		Int64("height", tx.Height).Msg("handled message")

	return nil
}

// -----------------------------------------------------------------------------------------------------

// handleMsgCreateRelationship allows to handle a MsgCreateRelationship properly
func (m *Module) handleMsgCreateRelationship(tx *juno.Tx, msg *relationshipstypes.MsgCreateRelationship) error {
	return m.db.SaveRelationship(types.NewRelationship(
		relationshipstypes.NewRelationship(msg.Signer, msg.Counterparty, msg.SubspaceID),
		tx.Height,
	))
}

// handleMsgDeleteRelationship allows to handle a MsgDeleteRelationship properly
func (m *Module) handleMsgDeleteRelationship(tx *juno.Tx, msg *relationshipstypes.MsgDeleteRelationship) error {
	return m.db.DeleteRelationship(types.NewRelationship(
		relationshipstypes.NewRelationship(msg.Signer, msg.Counterparty, msg.SubspaceID),
		tx.Height,
	))
}

// -----------------------------------------------------------------------------------------------------

// handleMsgBlockUser allows to handle a MsgBlockUser properly
func (m *Module) handleMsgBlockUser(tx *juno.Tx, msg *relationshipstypes.MsgBlockUser) error {
	return m.db.SaveUserBlock(types.NewBlockage(
		relationshipstypes.NewUserBlock(
			msg.Blocker,
			msg.Blocked,
			msg.Reason,
			msg.SubspaceID,
		),
		tx.Height,
	))
}

// handleMsgUnblockUser allows to handle a MsgUnblockUser properly
func (m *Module) handleMsgUnblockUser(tx *juno.Tx, msg *relationshipstypes.MsgUnblockUser) error {
	return m.db.DeleteBlockage(types.NewBlockage(
		relationshipstypes.NewUserBlock(msg.Blocker, msg.Blocked, "", msg.SubspaceID),
		tx.Height,
	))
}
