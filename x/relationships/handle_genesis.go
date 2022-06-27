package relationships

import (
	"encoding/json"

	relationshipstypes "github.com/desmos-labs/desmos/v4/x/relationships/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/desmos-labs/djuno/v2/types"
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
		err := m.db.SaveBlockage(types.NewBlockage(blockage, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	return nil
}
