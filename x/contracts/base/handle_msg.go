package contracts

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	juno "github.com/forbole/juno/v3/types"

	"github.com/desmos-labs/djuno/v2/types"
)

// HandleMsgInstantiateContract handles a MsgInstantiateContract instance by refreshing the stored tips contracts
func (m *Module) HandleMsgInstantiateContract(tx *juno.Tx, index int, _ *wasmtypes.MsgInstantiateContract, contractType string) error {
	address, err := m.ParseContractAddress(tx, index)
	if err != nil {
		return err
	}

	// Refresh all the contracts for the code id of the tips contract
	return m.db.SaveContract(types.NewContract(address, contractType, tx.Height))
}
