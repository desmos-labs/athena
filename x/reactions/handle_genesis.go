package reactions

import (
	"encoding/json"

	reactionstypes "github.com/desmos-labs/desmos/v4/x/reactions/types"

	"github.com/desmos-labs/djuno/v2/types"

	tmtypes "github.com/tendermint/tendermint/types"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	var genState reactionstypes.GenesisState
	m.cdc.MustUnmarshalJSON(appState[reactionstypes.ModuleName], &genState)

	// Save registered reactions
	for _, reaction := range genState.RegisteredReactions {
		err := m.db.SaveRegisteredReaction(types.NewRegisteredReaction(reaction, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	// Save reactions
	for _, reaction := range genState.Reactions {
		err := m.db.SaveReaction(types.NewReaction(reaction, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	// Save params
	for _, params := range genState.SubspacesParams {
		err := m.db.SaveReactionParams(types.NewReactionParams(params, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	return nil
}
