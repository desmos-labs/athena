package subspaces

import (
	"github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/desmos-labs/athena/v2/x/filters"

	"github.com/rs/zerolog/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v5/types"

	subspacestypes "github.com/desmos-labs/desmos/v6/x/subspaces/types"

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
	case *subspacestypes.MsgCreateSubspace:
		return m.handleMsgCreateSubspace(tx, index)

	case *subspacestypes.MsgEditSubspace:
		return m.handleMsgEditSubspace(tx, desmosMsg)

	case *subspacestypes.MsgDeleteSubspace:
		return m.handleMsgDeleteSubspace(tx, desmosMsg)

	case *subspacestypes.MsgCreateSection:
		return m.handleMsgCreateSection(tx, index, desmosMsg)

	case *subspacestypes.MsgEditSection:
		return m.handleMsgEditSection(tx, desmosMsg)

	case *subspacestypes.MsgMoveSection:
		return m.handleMsgMoveSection(tx, desmosMsg)

	case *subspacestypes.MsgDeleteSection:
		return m.handleMsgDeleteSection(tx, desmosMsg)

	case *subspacestypes.MsgCreateUserGroup:
		return m.handleMsgCreateUserGroup(tx, index, desmosMsg)

	case *subspacestypes.MsgEditUserGroup:
		return m.handleMsgEditUserGroup(tx, desmosMsg)

	case *subspacestypes.MsgMoveUserGroup:
		return m.handleMsgMoveUserGroup(tx, desmosMsg)

	case *subspacestypes.MsgSetUserGroupPermissions:
		return m.handleMsgSetUserGroupPermissions(tx, desmosMsg)

	case *subspacestypes.MsgDeleteUserGroup:
		return m.handleMsgDeleteUserGroup(tx, desmosMsg)

	case *subspacestypes.MsgAddUserToUserGroup:
		return m.handleMsgAddUserToUserGroup(tx, desmosMsg)

	case *subspacestypes.MsgRemoveUserFromUserGroup:
		return m.handleMsgRemoveUserFromUserGroup(tx, desmosMsg)

	case *subspacestypes.MsgSetUserPermissions:
		return m.handleMsgSetUserPermissions(tx, desmosMsg)
	}

	log.Debug().Str("module", "subspaces").Str("message", sdk.MsgTypeURL(msg)).
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

	return m.RefreshSubspaceData(tx.Height, subspaceID)
}

// handleMsgEditSubspace handles a MsgEditSubspace
func (m *Module) handleMsgEditSubspace(tx *juno.Tx, msg *subspacestypes.MsgEditSubspace) error {
	return m.updateSubspace(tx.Height, msg.SubspaceID)
}

// handleMsgDeleteSubspace handles a MsgDeleteSubspace
func (m *Module) handleMsgDeleteSubspace(tx *juno.Tx, msg *subspacestypes.MsgDeleteSubspace) error {
	return m.db.DeleteSubspace(tx.Height, msg.SubspaceID)
}

// -----------------------------------------------------------------------------------------------------

// handleMsgCreateSection handles a MsgCreateSection
func (m *Module) handleMsgCreateSection(tx *juno.Tx, index int, msg *subspacestypes.MsgCreateSection) error {
	// Get the subspace id
	event, err := tx.FindEventByType(index, subspacestypes.EventTypeCreateSection)
	if err != nil {
		return err
	}
	sectionIDStr, err := tx.FindAttributeByKey(event, subspacestypes.AttributeKeySectionID)
	if err != nil {
		return err
	}
	sectionID, err := subspacestypes.ParseSectionID(sectionIDStr)
	if err != nil {
		return err
	}

	return m.updateSection(tx.Height, msg.SubspaceID, sectionID)
}

// handleMsgEditSection handles a MsgEditSection
func (m *Module) handleMsgEditSection(tx *juno.Tx, msg *subspacestypes.MsgEditSection) error {
	return m.updateSection(tx.Height, msg.SubspaceID, msg.SectionID)
}

// handleMsgMoveSection handles a MsgMoveSection
func (m *Module) handleMsgMoveSection(tx *juno.Tx, msg *subspacestypes.MsgMoveSection) error {
	return m.updateSection(tx.Height, msg.SubspaceID, msg.SectionID)
}

// handleMsgDeleteSection handles a MsgDeleteSection
func (m *Module) handleMsgDeleteSection(tx *juno.Tx, msg *subspacestypes.MsgDeleteSection) error {
	return m.db.DeleteSection(tx.Height, msg.SubspaceID, msg.SectionID)
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

	// Update the user group
	err = m.updateUserGroup(tx.Height, msg.SubspaceID, groupID)
	if err != nil {
		return err
	}

	// Handle initial members
	for _, member := range msg.InitialMembers {
		err = m.db.AddUserToGroup(types.NewUserGroupMember(msg.SubspaceID, groupID, member, tx.Height))
		if err != nil {
			return err
		}
	}

	return nil
}

// handleMsgEditUserGroup handles a MsgEditUserGroup
func (m *Module) handleMsgEditUserGroup(tx *juno.Tx, msg *subspacestypes.MsgEditUserGroup) error {
	return m.updateUserGroup(tx.Height, msg.SubspaceID, msg.GroupID)
}

// handleMsgMoveUserGroup handles a MsgMoveUserGroup
func (m *Module) handleMsgMoveUserGroup(tx *juno.Tx, msg *subspacestypes.MsgMoveUserGroup) error {
	return m.updateUserGroup(tx.Height, msg.SubspaceID, msg.GroupID)
}

// handleMsgSetUserGroupPermissions handles a MsgSetUserGroupPermissions properly
func (m *Module) handleMsgSetUserGroupPermissions(tx *juno.Tx, msg *subspacestypes.MsgSetUserGroupPermissions) error {
	return m.updateUserGroup(tx.Height, msg.SubspaceID, msg.GroupID)
}

// handleMsgDeleteUserGroup handles a MsgDeleteUserGroup
func (m *Module) handleMsgDeleteUserGroup(tx *juno.Tx, msg *subspacestypes.MsgDeleteUserGroup) error {
	return m.db.DeleteUserGroup(tx.Height, msg.SubspaceID, msg.GroupID)
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
	return m.updateUserPermissions(tx.Height, msg.SubspaceID, msg.SectionID, msg.User)
}
