package types

import subspacestypes "github.com/desmos-labs/desmos/v3/x/subspaces/types"

type Subspace struct {
	subspacestypes.Subspace
	Height int64
}

func NewSubspace(subspace subspacestypes.Subspace, height int64) Subspace {
	return Subspace{
		Subspace: subspace,
		Height:   height,
	}
}

type UserPermission struct {
	subspacestypes.ACLEntry
	Height int64
}

func NewUserPermission(permission subspacestypes.ACLEntry, height int64) UserPermission {
	return UserPermission{
		ACLEntry: permission,
		Height:   height,
	}
}

type UserGroup struct {
	subspacestypes.UserGroup
	Height int64
}

func NewUserGroup(group subspacestypes.UserGroup, height int64) UserGroup {
	return UserGroup{
		UserGroup: group,
		Height:    height,
	}
}

type UserGroupMember struct {
	SubspaceID uint64
	GroupID    uint32
	Member     string
	Height     int64
}

func NewUserGroupMember(subspaceID uint64, groupID uint32, member string, height int64) UserGroupMember {
	return UserGroupMember{
		SubspaceID: subspaceID,
		GroupID:    groupID,
		Member:     member,
		Height:     height,
	}
}
