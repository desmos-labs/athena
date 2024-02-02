package relationships

import (
	"encoding/json"

	tmtypes "github.com/cometbft/cometbft/types"

	relationshipstypes "github.com/desmos-labs/desmos/v6/x/relationships/types"

	"github.com/desmos-labs/athena/v2/types"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	var genState relationshipstypes.GenesisState
	m.cdc.MustUnmarshalJSON(appState[relationshipstypes.ModuleName], &genState)

	// Save relationships
	for _, relationship := range genState.Relationships {
		err := m.db.SaveRelationship(types.NewRelationship(relationship, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	// Save blockages
	for _, blockage := range genState.Blocks {
		err := m.db.SaveUserBlock(types.NewBlockage(blockage, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	return nil
}
