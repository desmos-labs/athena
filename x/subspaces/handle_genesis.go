package subspaces

import (
	"encoding/json"

	"github.com/desmos-labs/djuno/v2/types"

	tmtypes "github.com/cometbft/cometbft/types"

	subspacestypes "github.com/desmos-labs/desmos/v6/x/subspaces/types"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	var genState subspacestypes.GenesisState
	m.cdc.MustUnmarshalJSON(appState[subspacestypes.ModuleName], &genState)

	// Save subspaces
	for _, subspace := range genState.Subspaces {
		err := m.db.SaveSubspace(types.NewSubspace(subspace, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	// Save sections
	for _, section := range genState.Sections {
		err := m.db.SaveSection(types.NewSection(section, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	// Save user permissions
	for _, permission := range genState.UserPermissions {
		err := m.db.SaveUserPermission(types.NewUserPermission(permission, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	// Save user groups
	for _, group := range genState.UserGroups {
		err := m.db.SaveUserGroup(types.NewUserGroup(group, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	// Save user group members
	for _, entry := range genState.UserGroupsMembers {
		err := m.db.AddUserToGroup(types.NewUserGroupMember(entry.SubspaceID, entry.GroupID, entry.User, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	return nil
}
