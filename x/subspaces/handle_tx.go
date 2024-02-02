package subspaces

import (
	abci "github.com/cometbft/cometbft/abci/types"
	subspacestypes "github.com/desmos-labs/desmos/v6/x/subspaces/types"
	juno "github.com/forbole/juno/v5/types"

	"github.com/desmos-labs/athena/types"
	"github.com/desmos-labs/athena/utils/transactions"
)

func (m *Module) HandleTx(tx *juno.Tx) error {
	return transactions.ParseTxEvents(tx, map[string]func(tx *juno.Tx, event abci.Event) error{
		subspacestypes.EventTypeCreateSubspace: m.handleCreateSubspaceEvent,
		subspacestypes.EventTypeEditSubspace:   m.handleEditSubspaceEvent,
		subspacestypes.EventTypeDeleteSubspace: m.handleDeleteSubspaceEvent,

		subspacestypes.EventTypeCreateSection: m.handleCreateSectionEvent,
		subspacestypes.EventTypeEditSection:   m.handleEditSectionEvent,
		subspacestypes.EventTypeMoveSection:   m.handleMoveSectionEvent,
		subspacestypes.EventTypeDeleteSection: m.handleDeleteSectionEvent,

		subspacestypes.EventTypeCreateUserGroup:         m.handleCreateUserGroupEvent,
		subspacestypes.EventTypeEditUserGroup:           m.handleEditUserGroupEvent,
		subspacestypes.EvenTypeMoveUserGroup:            m.handleMoveUserGroupEvent,
		subspacestypes.EventTypeSetUserGroupPermissions: m.handleSetUserGroupPermissionsEvent,
		subspacestypes.EventTypeDeleteUserGroup:         m.handleDeleteUserGroupEvent,

		subspacestypes.EventTypeAddUserToGroup:      m.handleAddUserToGroupEvent,
		subspacestypes.EventTypeRemoveUserFromGroup: m.handleRemoveUserFromGroupEvent,

		subspacestypes.EventTypeSetUserPermissions: m.handleSetUserPermissionsEvent,
	})
}

// --------------------------------------------------------------------------------------------------------------------

// handleCreateSubspaceEvent handles the creation of a new subspace
func (m *Module) handleCreateSubspaceEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	return m.RefreshSubspaceData(tx.Height, subspaceID)
}

// handleEditSubspaceEvent handles the edit of an existing subspace
func (m *Module) handleEditSubspaceEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	return m.updateSubspace(tx.Height, subspaceID)
}

// handleDeleteSubspaceEvent handles the deletion of an existing subspace
func (m *Module) handleDeleteSubspaceEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	return m.db.DeleteSubspace(tx.Height, subspaceID)
}

// --------------------------------------------------------------------------------------------------------------------

// handleCreateSectionEvent handles the creation of a new section
func (m *Module) handleCreateSectionEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	sectionID, err := GetSectionIDFromEvent(event)
	if err != nil {
		return err
	}

	return m.updateSection(tx.Height, subspaceID, sectionID)
}

// handleEditSectionEvent handles the edit of an existing section
func (m *Module) handleEditSectionEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	sectionID, err := GetSectionIDFromEvent(event)
	if err != nil {
		return err
	}

	return m.updateSection(tx.Height, subspaceID, sectionID)
}

func (m *Module) handleMoveSectionEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	sectionID, err := GetSectionIDFromEvent(event)
	if err != nil {
		return err
	}

	return m.updateSection(tx.Height, subspaceID, sectionID)
}

// handleDeleteSectionEvent handles the deletion of an existing section
func (m *Module) handleDeleteSectionEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	sectionID, err := GetSectionIDFromEvent(event)
	if err != nil {
		return err
	}

	return m.db.DeleteSection(tx.Height, subspaceID, sectionID)
}

// --------------------------------------------------------------------------------------------------------------------

// handleCreateUserGroupEvent handles the creation of a new user group
func (m *Module) handleCreateUserGroupEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}
	groupID, err := GetUserGroupIDFromEvent(event)
	if err != nil {
		return err
	}

	// Update the user group
	return m.updateUserGroup(tx.Height, subspaceID, groupID)
}

// handleEditUserGroupEvent handles the edit of an existing user group
func (m *Module) handleEditUserGroupEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}
	groupID, err := GetUserGroupIDFromEvent(event)
	if err != nil {
		return err
	}

	// Update the user group
	return m.updateUserGroup(tx.Height, subspaceID, groupID)
}

// handleMoveUserGroupEvent handles the move of an existing user group
func (m *Module) handleMoveUserGroupEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}
	groupID, err := GetUserGroupIDFromEvent(event)
	if err != nil {
		return err
	}

	// Update the user group
	return m.updateUserGroup(tx.Height, subspaceID, groupID)
}

// handleSetUserGroupPermissionsEvent handles the setting of permissions for an existing user group
func (m *Module) handleSetUserGroupPermissionsEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}
	groupID, err := GetUserGroupIDFromEvent(event)
	if err != nil {
		return err
	}

	// Update the user group
	return m.updateUserGroup(tx.Height, subspaceID, groupID)
}

// handleDeleteUserGroupEvent handles the deletion of an existing user group
func (m *Module) handleDeleteUserGroupEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}
	groupID, err := GetUserGroupIDFromEvent(event)
	if err != nil {
		return err
	}

	// Update the user group
	return m.db.DeleteUserGroup(tx.Height, subspaceID, groupID)
}

// --------------------------------------------------------------------------------------------------------------------

// handleAddUserToGroupEvent handles the addition of a user to an existing user group
func (m *Module) handleAddUserToGroupEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}
	groupID, err := GetUserGroupIDFromEvent(event)
	if err != nil {
		return err
	}
	member, err := GetUserFromEvent(event)
	if err != nil {
		return err
	}

	return m.db.AddUserToGroup(types.NewUserGroupMember(subspaceID, groupID, member, tx.Height))
}

// handleRemoveUserFromGroupEvent handles the removal of a user from an existing user group
func (m *Module) handleRemoveUserFromGroupEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}
	groupID, err := GetUserGroupIDFromEvent(event)
	if err != nil {
		return err
	}
	member, err := GetUserFromEvent(event)
	if err != nil {
		return err
	}

	return m.db.RemoveUserFromGroup(types.NewUserGroupMember(subspaceID, groupID, member, tx.Height))
}

// --------------------------------------------------------------------------------------------------------------------

// handleSetUserPermissionsEvent handles the setting of permissions for an existing user
func (m *Module) handleSetUserPermissionsEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}
	user, err := GetUserFromEvent(event)
	if err != nil {
		return err
	}

	return m.updateUserPermissions(tx.Height, subspaceID, user)
}
