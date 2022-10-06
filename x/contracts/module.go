package contracts

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"
	juno "github.com/forbole/juno/v3/types"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/forbole/juno/v3/node/remote"

	"github.com/desmos-labs/djuno/v2/database"
	"github.com/desmos-labs/djuno/v2/types"
)

// Module represents a generic smart contract module that can be extended for custom contracts handling
type Module struct {
	wasmClient wasmtypes.QueryClient
	db         *database.Db
}

// NewModule returns a new Module instance
func NewModule(wasmClient wasmtypes.QueryClient, db *database.Db) *Module {
	return &Module{
		wasmClient: wasmClient,
		db:         db,
	}
}

// HandleMsgInstantiateContract handles a MsgInstantiateContract instance by refreshing the stored tips contracts
func (m *Module) HandleMsgInstantiateContract(tx *juno.Tx, index int, _ *wasmtypes.MsgInstantiateContract, contractType string) error {
	event, err := tx.FindEventByType(index, wasmtypes.EventTypeInstantiate)
	if err != nil {
		return fmt.Errorf("no even with type %s found", wasmtypes.EventTypeInstantiate)
	}
	address, err := tx.FindAttributeByKey(event, wasmtypes.AttributeKeyContractAddr)
	if err != nil {
		return fmt.Errorf("no %s attribute found", wasmtypes.AttributeKeyContractAddr)
	}

	// Refresh all the contracts for the code id of the tips contract
	return m.db.SaveContract(types.NewContract(address, contractType, tx.Height))
}

// RefreshContracts refreshes the contracts that have been instanciated for the given code id, at the given height.
// After fetching the data from the chain, such contracts addresses are stored
// inside the database as contracts of the given type
func (m *Module) RefreshContracts(height int64, codeID uint64, contractType string) error {
	contracts, err := m.getAllContractsByCode(height, codeID)
	if err != nil {
		return fmt.Errorf("error while getting contracts for code id %d: %s", codeID, err)
	}

	for _, contract := range contracts {
		err = m.db.SaveContract(types.NewContract(contract, contractType, height))
		if err != nil {
			return err
		}
	}

	return nil
}

// getAllContractsByCode returns all the contracts addresses having the given code id at the given height
func (m *Module) getAllContractsByCode(height int64, codeID uint64) ([]string, error) {
	var contracts []string
	var nextKey []byte
	var stop bool

	for !stop {
		res, err := m.wasmClient.ContractsByCode(
			remote.GetHeightRequestContext(context.Background(), height),
			&wasmtypes.QueryContractsByCodeRequest{
				CodeId: codeID,
				Pagination: &query.PageRequest{
					Key: nextKey,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		contracts = append(contracts, res.Contracts...)
		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return contracts, nil
}
