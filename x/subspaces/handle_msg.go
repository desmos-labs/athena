package subspaces

import (
	"github.com/gogo/protobuf/proto"

	"github.com/rs/zerolog/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	subspacestypes "github.com/desmos-labs/desmos/v3/x/subspaces/types"
	juno "github.com/forbole/juno/v3/types"

	"github.com/desmos-labs/djuno/v2/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch desmosMsg := msg.(type) {
	case *subspacestypes.MsgCreateSubspace:
		return m.handleMsgCreateSubspace(tx, index)

	case *subspacestypes.MsgEditSubspace:
		return m.handleMsgEditSubspace(tx, desmosMsg)

	case *subspacestypes.MsgDeleteSubspace:
		return m.handleMsgDeleteSubspace(desmosMsg)

	case *subspacestypes.MsgCreateUserGroup:
		return m.handleMsgCreateUserGroup(tx, index, desmosMsg)

	case *subspacestypes.MsgEditUserGroup:
		return m.handleMsgEditUserGroup(tx, desmosMsg)

	case *subspacestypes.MsgSetUserGroupPermissions:
		return m.handleMsgSetUserGroupPermissions(tx, desmosMsg)

	case *subspacestypes.MsgDeleteUserGroup:
		return m.handleMsgDeleteUserGroup(desmosMsg)

	case *subspacestypes.MsgAddUserToUserGroup:
		return m.handleMsgAddUserToUserGroup(tx, desmosMsg)

	case *subspacestypes.MsgRemoveUserFromUserGroup:
		return m.handleMsgRemoveUserFromUserGroup(tx, desmosMsg)

	case *subspacestypes.MsgSetUserPermissions:
		return m.handleMsgSetUserPermissions(tx, desmosMsg)
	}

	log.Debug().Str("module", "subspaces").Str("message", proto.MessageName(msg)).
		Int64("height", tx.Height).Msg("handled message")

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

// handleMsgCreateSubspace handles a MsgCreateSubspace
func (m *Module) handleMsgCreateSubspace(tx *juno.Tx, index int) error {
	// Get the subspace id
	event, err := tx.FindEventByType(index, subspacestypes.EventTypeCreateSubspace)
	if err != nil {
		return err
	}
	subspaceIDStr, err := tx.FindAttributeByKey(event, subspacestypes.AttributeKeySubspaceID)
	if err != nil {
		return err
	}
	subspaceID, err := subspacestypes.ParseSubspaceID(subspaceIDStr)
	if err != nil {
		return err
	}

	return m.updateSubspace(tx.Height, subspaceID)
}

// handleMsgEditSubspace handles a MsgEditSubspace
func (m *Module) handleMsgEditSubspace(tx *juno.Tx, msg *subspacestypes.MsgEditSubspace) error {
	return m.updateSubspace(tx.Height, msg.SubspaceID)
}

// handleMsgDeleteSubspace handles a MsgDeleteSubspace
func (m *Module) handleMsgDeleteSubspace(msg *subspacestypes.MsgDeleteSubspace) error {
	return m.db.DeleteSubspace(msg.SubspaceID)
}

// -----------------------------------------------------------------------------------------------------

// handleMsgCreateUserGroup handles a MsgCreateUserGroup
func (m *Module) handleMsgCreateUserGroup(tx *juno.Tx, index int, msg *subspacestypes.MsgCreateUserGroup) error {
	// Get the group id
	event, err := tx.FindEventByType(index, subspacestypes.EventTypeCreateUserGroup)
	if err != nil {
		return err
	}
	groupIDStr, err := tx.FindAttributeByKey(event, subspacestypes.AttributeKeyUserGroupID)
	if err != nil {
		return err
	}
	groupID, err := subspacestypes.ParseGroupID(groupIDStr)
	if err != nil {
		return err
	}

	return m.updateUserGroup(tx.Height, msg.SubspaceID, groupID)
}

// handleMsgEditUserGroup handles a MsgEditUserGroup
func (m *Module) handleMsgEditUserGroup(tx *juno.Tx, msg *subspacestypes.MsgEditUserGroup) error {
	return m.updateUserGroup(tx.Height, msg.SubspaceID, msg.GroupID)
}

// handleMsgSetUserGroupPermissions handles a MsgSetUserGroupPermissions properly
func (m *Module) handleMsgSetUserGroupPermissions(tx *juno.Tx, msg *subspacestypes.MsgSetUserGroupPermissions) error {
	return m.updateUserGroup(tx.Height, msg.SubspaceID, msg.GroupID)
}

// handleMsgDeleteUserGroup handles a MsgDeleteUserGroup
func (m *Module) handleMsgDeleteUserGroup(msg *subspacestypes.MsgDeleteUserGroup) error {
	return m.db.DeleteUserGroup(msg.SubspaceID, msg.GroupID)
}

// -----------------------------------------------------------------------------------------------------

// handleMsgAddUserToUserGroup handles a MsgAddUserToUserGroup
func (m *Module) handleMsgAddUserToUserGroup(tx *juno.Tx, msg *subspacestypes.MsgAddUserToUserGroup) error {
	return m.db.AddUserToGroup(types.NewUserGroupMember(msg.SubspaceID, msg.GroupID, msg.User, tx.Height))
}

// handleMsgRemoveUserFromUserGroup handles a MsgRemoveUserFromUserGroup
func (m *Module) handleMsgRemoveUserFromUserGroup(tx *juno.Tx, msg *subspacestypes.MsgRemoveUserFromUserGroup) error {
	return m.db.RemoveUserFromGroup(types.NewUserGroupMember(msg.SubspaceID, msg.GroupID, msg.User, tx.Height))
}

// -----------------------------------------------------------------------------------------------------

// handleMsgSetUserPermissions handles a MsgSetUserPermissions
func (m *Module) handleMsgSetUserPermissions(tx *juno.Tx, msg *subspacestypes.MsgSetUserPermissions) error {
	return m.updateUserPermissions(tx.Height, msg.SubspaceID, msg.User)
}
