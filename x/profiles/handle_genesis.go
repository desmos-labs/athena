package profiles

import (
	"encoding/json"

	tmtypes "github.com/cometbft/cometbft/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	profilestypes "github.com/desmos-labs/desmos/v7/x/profiles/types"

	"github.com/desmos-labs/athena/v2/types"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	var authGenState authtypes.GenesisState
	m.cdc.MustUnmarshalJSON(appState[authtypes.ModuleName], &authGenState)

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

		err = m.db.SaveProfile(types.NewProfile(profile, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	var genState profilestypes.GenesisState
	m.cdc.MustUnmarshalJSON(appState[profilestypes.ModuleName], &genState)

	// Save DTag transfer requests
	for _, request := range genState.DTagTransferRequests {
		err = m.db.SaveDTagTransferRequest(types.NewDTagTransferRequest(request, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	// Save chain links
	for _, link := range genState.ChainLinks {
		err = m.db.SaveChainLink(types.NewChainLink(link, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	// Save application links
	for _, link := range genState.ApplicationLinks {
		err = m.db.SaveApplicationLink(types.NewApplicationLink(link, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	// Save params
	err = m.db.SaveProfilesParams(types.NewProfilesParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return err
	}

	return nil
}
