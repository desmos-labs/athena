package subspaces

import (
	"context"

	"github.com/forbole/juno/v4/node/remote"

	subspacestypes "github.com/desmos-labs/desmos/v4/x/subspaces/types"

	"github.com/desmos-labs/djuno/v2/types"
)

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
func (m *Module) updateUserPermissions(height int64, subspaceID uint64, sectionID uint32, user string) error {
	// Get the permissions
	res, err := m.client.UserPermissions(
		remote.GetHeightRequestContext(context.Background(), height),
		&subspacestypes.QueryUserPermissionsRequest{
			SubspaceId: subspaceID,
			SectionId:  sectionID,
			User:       user,
		},
	)
	if err != nil {
		return err
	}

	// Save the user permissions
	return m.db.SaveUserPermission(types.NewUserPermission(
		subspacestypes.NewUserPermission(subspaceID, sectionID, user, res.Permissions),
		height,
	))
}
