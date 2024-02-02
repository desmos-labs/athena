package subspaces

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/forbole/juno/v5/types/utils"

	"github.com/forbole/juno/v5/node/remote"

	subspacestypes "github.com/desmos-labs/desmos/v6/x/subspaces/types"

	"github.com/desmos-labs/athena/types"
)

// GetSubspaceIDFromEvent returns the subspace ID from the given event
func GetSubspaceIDFromEvent(event abci.Event) (uint64, error) {
	subspaceIDStr, err := utils.FindAttributeByKey(event, subspacestypes.AttributeKeySubspaceID)
	if err != nil {
		return 0, err
	}
	return subspacestypes.ParseSubspaceID(subspaceIDStr.Value)
}

// GetSectionIDFromEvent returns the section ID from the given event
func GetSectionIDFromEvent(event abci.Event) (uint32, error) {
	sectionIDStr, err := utils.FindAttributeByKey(event, subspacestypes.AttributeKeySectionID)
	if err != nil {
		return 0, err
	}
	return subspacestypes.ParseSectionID(sectionIDStr.Value)
}

// GetUserGroupIDFromEvent returns the user group ID from the given event
func GetUserGroupIDFromEvent(event abci.Event) (uint32, error) {
	groupIDStr, err := utils.FindAttributeByKey(event, subspacestypes.AttributeKeyUserGroupID)
	if err != nil {
		return 0, err
	}
	return subspacestypes.ParseGroupID(groupIDStr.Value)
}

// GetUserFromEvent returns the user from the given event
func GetUserFromEvent(event abci.Event) (string, error) {
	userStr, err := utils.FindAttributeByKey(event, subspacestypes.AttributeKeyUser)
	if err != nil {
		return "", err
	}
	return userStr.Value, nil
}

// --------------------------------------------------------------------------------------------------------------------

// updateSubspace updates the stored data for the given subspace at the specified height
func (m *Module) updateSubspace(height int64, subspaceID uint64) error {
	// Get the subspace
	res, err := m.client.Subspace(
		remote.GetHeightRequestContext(context.Background(), height),
		subspacestypes.NewQuerySubspaceRequest(subspaceID),
	)
	if err != nil {
		return err
	}

	// Save the subspace
	return m.db.SaveSubspace(types.NewSubspace(res.Subspace, height))
}

// updateSection updates the stored data for the given subspace section at the specified height
func (m *Module) updateSection(height int64, subspaceID uint64, sectionID uint32) error {
	// Get the subspace
	res, err := m.client.Section(
		remote.GetHeightRequestContext(context.Background(), height),
		subspacestypes.NewQuerySectionRequest(subspaceID, sectionID),
	)
	if err != nil {
		return err
	}

	// Save the subspace
	return m.db.SaveSection(types.NewSection(res.Section, height))
}

// updateUserGroup updates the stored data for the given user group at the specified height
func (m *Module) updateUserGroup(height int64, subspaceID uint64, groupID uint32) error {
	// Get the user group
	res, err := m.client.UserGroup(
		remote.GetHeightRequestContext(context.Background(), height),
		subspacestypes.NewQueryUserGroupRequest(subspaceID, groupID),
	)
	if err != nil {
		return err
	}

	// Save the user group
	return m.db.SaveUserGroup(types.NewUserGroup(res.Group, height))
}

// updateUserPermissions updates the stored permissions for the given user at the specified height
func (m *Module) updateUserPermissions(height int64, subspaceID uint64, user string) error {
	// Get the sections
	sectionsRes, err := m.client.Sections(
		remote.GetHeightRequestContext(context.Background(), height),
		&subspacestypes.QuerySectionsRequest{
			SubspaceId: subspaceID,
		},
	)
	if err != nil {
		return err
	}

	for _, section := range sectionsRes.Sections {
		// Get the permissions
		res, err := m.client.UserPermissions(
			remote.GetHeightRequestContext(context.Background(), height),
			&subspacestypes.QueryUserPermissionsRequest{
				SubspaceId: subspaceID,
				SectionId:  section.ID,
				User:       user,
			},
		)
		if err != nil {
			return err
		}

		// Save the user permissions
		err = m.db.SaveUserPermission(types.NewUserPermission(
			subspacestypes.NewUserPermission(subspaceID, section.ID, user, res.Permissions),
			height,
		))
		if err != nil {
			return err
		}
	}

	return nil
}
