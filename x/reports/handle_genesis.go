package reports

import (
	"encoding/json"

	"github.com/desmos-labs/athena/v2/types"

	tmtypes "github.com/cometbft/cometbft/types"

	reportstypes "github.com/desmos-labs/desmos/v7/x/reports/types"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	var genState reportstypes.GenesisState
	m.cdc.MustUnmarshalJSON(appState[reportstypes.ModuleName], &genState)

	// Save reasons
	for _, subspace := range genState.Reasons {
		err := m.db.SaveReason(types.NewReason(subspace, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	// Save reports
	for _, section := range genState.Reports {
		err := m.db.SaveReport(types.NewReport(section, doc.InitialHeight))
		if err != nil {
			return err
		}
	}

	// Save the params
	return m.db.SaveReportsParams(types.NewReportsParams(genState.Params, doc.InitialHeight))
}
