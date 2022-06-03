package database

import (
	dbtypes "github.com/desmos-labs/djuno/v2/database/types"
	"github.com/desmos-labs/djuno/v2/types"
)

// SaveSubspace stores the given subspace inside the database
func (db *Db) SaveSubspace(subspace types.Subspace) error {
	stmt := `
INSERT INTO subspace (id, name, description, treasury, owner, creator, creation_time, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT DO UPDATE 
    SET name = excluded.name,
        description = excluded.description,
        treasury = excluded.treasury,
        owner = excluded.owner,
        creator = excluded.creator,
        creation_time = excluded.creation_time,
        height = excluded.height
WHERE subspace.height <= excluded.height`

	_, err := db.Sql.Exec(stmt,
		subspace.ID,
		subspace.Name,
		dbtypes.ToNullString(subspace.Description),
		dbtypes.ToNullString(subspace.Treasury),
		subspace.Owner,
		subspace.Creator,
		subspace.CreationTime,
		subspace.Height,
	)
	return err
}

// DeleteSubspace removes the subspace with the given id from the database
func (db *Db) DeleteSubspace(id uint64) error {
	stmt := `DELETE FROM subspace WHERE id = $1`
	_, err := db.Sql.Exec(stmt, id)
	return err
}

// --------------------------------------------------------------------------------------------------------------------

// getUserGroupRowID returns the row id associated to the group with the given details
func (db *Db) getUserGroupRowID(subspaceID uint64, groupID uint32) (uint64, error) {
	stmt := `SELECT row_id FROM subspace_user_group WHERE subspace_id = $1 and id = $2`

	var rowID uint64
	err := db.Sql.QueryRow(stmt, subspaceID, groupID).Scan(&rowID)
	return rowID, err
}

// SaveUserGroup stores the given group inside the database
func (db *Db) SaveUserGroup(group types.UserGroup) error {
	stmt := `
INSERT INTO subspace_user_group (subspace_id, id, name, description, permissions, height) 
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT ON CONSTRAINT unique_subspace_user_group DO UPDATE 
    SET name = excluded.name,
        description = excluded.description,
        permissions = excluded.permissions,
        height = excluded.height
WHERE subspace_user_group.height <= excluded.height`

	_, err := db.Sql.Exec(stmt,
		group.SubspaceID,
		group.ID,
		group.Name,
		dbtypes.ToNullString(group.Description),
		group.Permissions,
		group.Height,
	)
	return err
}

// DeleteUserGroup removes the given user group from the subspace
func (db *Db) DeleteUserGroup(subspaceID uint64, groupID uint32) error {
	stmt := `DELETE FROM subspace_user_group WHERE subspace_id = $1 AND id = $2`
	_, err := db.Sql.Exec(stmt, subspaceID, groupID)
	return err
}

// AddUserToGroup adds a user to a user group
func (db *Db) AddUserToGroup(member types.UserGroupMember) error {
	rowID, err := db.getUserGroupRowID(member.SubspaceID, member.GroupID)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO subspace_user_group_member (group_row_id, member, height) 
VALUES ($1, $2, $3)
ON CONFLICT ON CONSTRAINT unique_subspace_group_membership DO NOTHING`

	_, err = db.Sql.Exec(stmt, rowID, member.Member, member.Height)
	return err
}

// RemoveUserFromGroup removes the member from the given user group
func (db *Db) RemoveUserFromGroup(member types.UserGroupMember) error {
	rowID, err := db.getUserGroupRowID(member.SubspaceID, member.GroupID)
	if err != nil {
		return err
	}

	stmt := `DELETE FROM subspace_user_group_member WHERE group_row_id = $1 AND member = $2 AND height <= $3`
	_, err = db.Sql.Exec(stmt, rowID, member.Member, member.Height)
	return err
}

// --------------------------------------------------------------------------------------------------------------------

// SaveUserPermission stores the given permissions inside the database
func (db *Db) SaveUserPermission(permission types.UserPermission) error {
	stmt := `
INSERT INTO subspace_user_permission (subspace_id, user_address, permissions, height) 
VALUES ($1, $2, $3, $4)
ON CONFLICT ON CONSTRAINT unique_subspace_permission DO UPDATE
    SET permissions = excluded.permissions,
        height = excluded.height
WHERE subspace_user_permission.height <= excluded.height`

	_, err := db.Sql.Exec(stmt, permission.SubspaceID, permission.User, permission.Permissions, permission.Height)
	return err
}

// --------------------------------------------------------------------------------------------------------------------
