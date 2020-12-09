package profiles

import (
	"encoding/json"
	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"

	"github.com/cosmos/cosmos-sdk/codec"
	desmosdb "github.com/desmos-labs/djuno/database"
)

// HandleGenesis allows to properly handle the genesis state for the posts module
func HandleGenesis(cdc *codec.LegacyAmino, appState map[string]json.RawMessage, db *desmosdb.DesmosDb) error {
	var genState profilestypes.GenesisState
	cdc.MustUnmarshalJSON(appState[profilestypes.ModuleName], &genState)

	// Save the profiles
	for _, prof := range genState.Profiles {
		err := db.SaveProfile(prof)
		if err != nil {
			return err
		}
	}

	return nil
}
