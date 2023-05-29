package types

import (
	subspacestypes "github.com/desmos-labs/desmos/v5/x/subspaces/types"
)

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

type Section struct {
	subspacestypes.Section
	Height int64
}

func NewSection(section subspacestypes.Section, height int64) Section {
	return Section{
		Section: section,
		Height:  height,
	}
}

type UserPermission struct {
	subspacestypes.UserPermission
	Height int64
}

func NewUserPermission(permission subspacestypes.UserPermission, height int64) UserPermission {
	return UserPermission{
		UserPermission: permission,
		Height:         height,
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
