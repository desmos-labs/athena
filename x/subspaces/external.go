package subspaces

import (
	"context"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/cosmos/cosmos-sdk/types/query"
	subspacestypes "github.com/desmos-labs/desmos/v4/x/subspaces/types"
	"github.com/forbole/juno/v3/node/remote"

	"github.com/desmos-labs/djuno/v2/types"
	"github.com/desmos-labs/djuno/v2/utils"
)

// RefreshSubspacesData refreshes all the subspaces user data
func (m *Module) RefreshSubspacesData(height int64) error {
	subspaces, err := m.QueryAllSubspaces(height)
	if err != nil {
		return err
	}

	err = m.db.DeleteAllSubspaces(height)
	if err != nil {
		return err
	}

	for _, subspace := range subspaces {
		// Save the subspace
		err = m.db.SaveSubspace(subspace)
		if err != nil {
			return err
		}

		// Update the sections
		sections, err := m.queryAllSections(height, subspace.ID)
		if err != nil {
			return err
		}

		for _, section := range sections {
			err = m.db.SaveSection(section)
			if err != nil {
				return err
			}
		}

		// Update the user groups
		groups, err := m.queryAllUserGroups(height, subspace.ID)
		if err != nil {
			return err
		}
		for _, group := range groups {
			err = m.db.SaveUserGroup(group)
			if err != nil {
				return err
			}

			// Update the members
			members, err := m.queryAllUserGroupMembers(height, group.SubspaceID, group.ID)
			if err != nil {
				return err
			}

			err = m.db.RemoveAllUsersFromGroup(height, group.SubspaceID, group.ID)
			if err != nil {
				return err
			}

			// Save the members
			for _, member := range members {
				err = m.db.AddUserToGroup(member)
				if err != nil {
					return err
				}
			}
		}

		// Update the user permissions
		permissions, err := m.queryAllUserPermissions(height, subspace.ID)
		if err != nil {
			return err
		}

		for _, permission := range permissions {
			err = m.db.SaveUserPermission(permission)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// QueryAllSubspaces queries all the subspaces present on the node at the given height
func (m *Module) QueryAllSubspaces(height int64) ([]types.Subspace, error) {
	var subspaces []types.Subspace

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.Subspaces(
			remote.GetHeightRequestContext(context.Background(), height),
			&subspacestypes.QuerySubspacesRequest{
				Pagination: &query.PageRequest{
					Key: nextKey,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, subspace := range res.Subspaces {
			subspaces = append(subspaces, types.NewSubspace(subspace, height))
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return subspaces, nil
}

// queryAllSections queries all the sections for the given subspace present on the node at the given height
func (m *Module) queryAllSections(height int64, subspaceID uint64) ([]types.Section, error) {
	var sections []types.Section

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.Sections(
			remote.GetHeightRequestContext(context.Background(), height),
			&subspacestypes.QuerySectionsRequest{
				SubspaceId: subspaceID,
				Pagination: &query.PageRequest{
					Key: nextKey,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, section := range res.Sections {
			sections = append(sections, types.NewSection(section, height))
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return sections, nil
}

// queryAllUserGroups queries all the user groups for the given subspace present on the node at the given height
func (m *Module) queryAllUserGroups(height int64, subspaceID uint64) ([]types.UserGroup, error) {
	var groups []types.UserGroup

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.UserGroups(
			remote.GetHeightRequestContext(context.Background(), height),
			&subspacestypes.QueryUserGroupsRequest{
				SubspaceId: subspaceID,
				Pagination: &query.PageRequest{
					Key: nextKey,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, group := range res.Groups {
			groups = append(groups, types.NewUserGroup(group, height))
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return groups, nil
}

// queryAllUserGroupMembers queries all the user group members for the given group present on the node at the given height
func (m *Module) queryAllUserGroupMembers(height int64, subspaceID uint64, groupID uint32) ([]types.UserGroupMember, error) {
	var members []types.UserGroupMember

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.UserGroupMembers(
			remote.GetHeightRequestContext(context.Background(), height),
			&subspacestypes.QueryUserGroupMembersRequest{
				SubspaceId: subspaceID,
				GroupId:    groupID,
				Pagination: &query.PageRequest{
					Key: nextKey,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, member := range res.Members {
			members = append(members, types.NewUserGroupMember(subspaceID, groupID, member, height))
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return members, nil
}

// queryAllUserPermissions queries all the user permissions for the given subspace present on the node at the given height
func (m *Module) queryAllUserPermissions(height int64, subspaceID uint64) ([]types.UserPermission, error) {
	// Query all the transactions
	permissionsQuery := fmt.Sprintf("%s.%s='%d' AND tx.height <= %d",
		subspacestypes.EventTypeSetUserPermissions,
		subspacestypes.AttributeKeySubspaceID,
		subspaceID,
		height,
	)
	txs, err := utils.QueryTxs(m.node, permissionsQuery)
	if err != nil {
		return nil, err
	}

	// Sort the txs based on their ascending height
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Height < txs[j].Height
	})

	// Parse all the transactions' messages
	var permissions []types.UserPermission
	for _, tx := range txs {
		transaction, err := m.node.Tx(hex.EncodeToString(tx.Tx.Hash()))
		if err != nil {
			return nil, err
		}

		// Handle only the MsgSetUserPermissions
		for _, msg := range transaction.GetMsgs() {
			if msg, ok := msg.(*subspacestypes.MsgSetUserPermissions); ok {
				permission := subspacestypes.NewUserPermission(msg.SubspaceID, msg.SectionID, msg.User, msg.Permissions)
				permissions = append(permissions, types.NewUserPermission(permission, height))
			}
		}
	}

	return permissions, nil
}
