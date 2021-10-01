package profiles

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/desmos-labs/djuno/types"

	desmosdb "github.com/desmos-labs/djuno/database"
)

// HandleGenesis allows to properly handle the genesis state for the posts module
func HandleGenesis(
	doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, cdc codec.Codec, db *desmosdb.Db,
) error {
	var authGenState authtypes.GenesisState
	cdc.MustUnmarshalJSON(appState[authtypes.ModuleName], &authGenState)

	accounts, err := authtypes.UnpackAccounts(authGenState.Accounts)
	if err != nil {
		return err
	}

	// Store the profiles
	for _, account := range accounts {
		profile, ok := account.(*profilestypes.Profile)
		if !ok {
			continue
		}

		err = db.SaveProfile(types.NewProfile(profile, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	var genState profilestypes.GenesisState
	cdc.MustUnmarshalJSON(appState[profilestypes.ModuleName], &genState)

	// Save DTag transfer requests
	for _, request := range genState.DTagTransferRequests {
		err = db.SaveDTagTransferRequest(types.NewDTagTransferRequest(request, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	// Save relationships
	for _, relationship := range genState.Relationships {
		err = db.SaveRelationship(types.NewRelationship(relationship, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	// Save blockages
	for _, blockage := range genState.Blocks {
		err = db.SaveBlockage(types.NewBlockage(blockage, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	// Save chain links
	for _, link := range genState.ChainLinks {
		err = db.SaveChainLink(types.NewChainLink(link, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	// Save application links
	for _, link := range genState.ApplicationLinks {
		err = db.SaveApplicationLink(types.NewApplicationLink(link, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	return nil
}
