package database

import (
	"database/sql"
	"errors"
	"fmt"

	dbtypes "github.com/desmos-labs/djuno/v2/database/types"
	"github.com/desmos-labs/djuno/v2/types"
)

// SaveSubspace stores the given subspace inside the database
func (db *Db) SaveSubspace(subspace types.Subspace) error {
	stmt := `
INSERT INTO subspace (id, name, description, treasury_address, owner_address, creator_address, creation_time, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (id) DO UPDATE 
    SET name = excluded.name,
        description = excluded.description,
        treasury_address = excluded.treasury_address,
        owner_address = excluded.owner_address,
        creator_address = excluded.creator_address,
        creation_time = excluded.creation_time,
        height = excluded.height
WHERE subspace.height <= excluded.height`

	_, err := db.SQL.Exec(stmt,
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
func (db *Db) DeleteSubspace(height int64, id uint64) error {
	stmt := `DELETE FROM subspace WHERE id = $1 AND height <= $2`
	_, err := db.SQL.Exec(stmt, id, height)
	return err
}

// DeleteAllSubspaces removes all the subspaces from the database
func (db *Db) DeleteAllSubspaces(height int64) error {
	stmt := `DELETE FROM subspace WHERE height <= $1`
	_, err := db.SQL.Exec(stmt, height)
	return err
}

// --------------------------------------------------------------------------------------------------------------------

// getUserGroupRowID returns the row id associated to the section with the given details
func (db *Db) getSectionRowID(subspaceID uint64, sectionID uint32) (sql.NullInt64, error) {
	stmt := `SELECT row_id FROM subspace_section WHERE subspace_id = $1 and id = $2`

	var rowID int64
	err := db.SQL.QueryRow(stmt, subspaceID, sectionID).Scan(&rowID)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.NullInt64{Int64: 0, Valid: false}, nil
	}

	return sql.NullInt64{Int64: rowID, Valid: true}, err
}

// SaveSection stores the given section inside the database
func (db *Db) SaveSection(section types.Section) error {
	parentRowID, err := db.getSectionRowID(section.SubspaceID, section.ParentID)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO subspace_section (subspace_id, id, parent_row_id, name, description, height) 
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT ON CONSTRAINT unique_subspace_section DO UPDATE 
    SET name = excluded.name,
        description = excluded.description,
        parent_row_id = excluded.parent_row_id,
        height = excluded.height
WHERE subspace_section.height <= excluded.height`

	_, err = db.SQL.Exec(stmt,
		section.SubspaceID,
		section.ID,
		parentRowID,
		section.Name,
		section.Description,
		section.Height,
	)
	return err
}

// DeleteSection removes the given section from the subspace
func (db *Db) DeleteSection(height int64, subspaceID uint64, sectionID uint32) error {
	stmt := `DELETE FROM subspace_section WHERE subspace_id = $1 AND id = $2 AND height <= $3`
	_, err := db.SQL.Exec(stmt, subspaceID, sectionID, height)
	return err
}

// --------------------------------------------------------------------------------------------------------------------

// getUserGroupRowID returns the row id associated to the group with the given details
func (db *Db) getUserGroupRowID(subspaceID uint64, groupID uint32) (uint64, error) {
	stmt := `SELECT row_id FROM subspace_user_group WHERE subspace_id = $1 and id = $2`

	var rowID uint64
	err := db.SQL.QueryRow(stmt, subspaceID, groupID).Scan(&rowID)
	return rowID, err
}

// SaveUserGroup stores the given group inside the database
func (db *Db) SaveUserGroup(group types.UserGroup) error {
	sectionRowID, err := db.getSectionRowID(group.SubspaceID, group.SectionID)
	if err != nil {
		return err
	}

	if !sectionRowID.Valid {
		return fmt.Errorf("section with id %d not found inside subspace %d", group.SectionID, group.SubspaceID)
	}

	stmt := `
INSERT INTO subspace_user_group (subspace_id, section_row_id, id, name, description, permissions, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT ON CONSTRAINT unique_subspace_user_group DO UPDATE 
    SET name = excluded.name,
        description = excluded.description,
        permissions = excluded.permissions,
        height = excluded.height
WHERE subspace_user_group.height <= excluded.height`

	_, err = db.SQL.Exec(stmt,
		group.SubspaceID,
		sectionRowID,
		group.ID,
		group.Name,
		dbtypes.ToNullString(group.Description),
		dbtypes.ConvertPermissions(group.Permissions),
		group.Height,
	)
	return err
}

// DeleteUserGroup removes the given user group from the subspace
func (db *Db) DeleteUserGroup(height int64, subspaceID uint64, groupID uint32) error {
	stmt := `DELETE FROM subspace_user_group WHERE subspace_id = $1 AND id = $2 AND height <= $3`
	_, err := db.SQL.Exec(stmt, subspaceID, groupID, height)
	return err
}

// AddUserToGroup adds a user to a user group
func (db *Db) AddUserToGroup(member types.UserGroupMember) error {
	rowID, err := db.getUserGroupRowID(member.SubspaceID, member.GroupID)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO subspace_user_group_member (group_row_id, member_address, height) 
VALUES ($1, $2, $3)
ON CONFLICT ON CONSTRAINT unique_subspace_group_membership DO NOTHING`

	_, err = db.SQL.Exec(stmt, rowID, member.Member, member.Height)
	return err
}

// RemoveUserFromGroup removes the member from the given user group
func (db *Db) RemoveUserFromGroup(member types.UserGroupMember) error {
	rowID, err := db.getUserGroupRowID(member.SubspaceID, member.GroupID)
	if err != nil {
		return err
	}

	stmt := `DELETE FROM subspace_user_group_member WHERE group_row_id = $1 AND member_address = $2 AND height <= $3`
	_, err = db.SQL.Exec(stmt, rowID, member.Member, member.Height)
	return err
}

// --------------------------------------------------------------------------------------------------------------------

// SaveUserPermission stores the given permissions inside the database
func (db *Db) SaveUserPermission(permission types.UserPermission) error {
	// TODO: Investigate why this happened
	if permission.Permissions == nil {
		return nil
	}

	sectionRowID, err := db.getSectionRowID(permission.SubspaceID, permission.SectionID)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO subspace_user_permission (section_row_id, user_address, permissions, height) 
VALUES ($1, $2, $3, $4)
ON CONFLICT ON CONSTRAINT unique_subspace_permission DO UPDATE
    SET permissions = excluded.permissions,
        height = excluded.height
WHERE subspace_user_permission.height <= excluded.height`

	_, err = db.SQL.Exec(stmt,
		sectionRowID,
		permission.User,
		dbtypes.ConvertPermissions(permission.Permissions),
		permission.Height,
	)
	return err
}

// --------------------------------------------------------------------------------------------------------------------
