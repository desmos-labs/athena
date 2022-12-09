package subspaces

import (
	"github.com/desmos-labs/djuno/v2/types"
)

type Database interface {
	SaveSubspace(subspace types.Subspace) error
	DeleteSubspace(height int64, id uint64) error
	DeleteAllSubspaces(height int64) error
	SaveSection(section types.Section) error
	DeleteSection(height int64, subspaceID uint64, sectionID uint32) error
	SaveUserGroup(group types.UserGroup) error
	DeleteUserGroup(height int64, subspaceID uint64, groupID uint32) error
	AddUserToGroup(member types.UserGroupMember) error
	RemoveUserFromGroup(member types.UserGroupMember) error
	SaveUserPermission(permission types.UserPermission) error
	DeleteUserPermission(permission types.UserPermission) error
}
